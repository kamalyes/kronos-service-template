/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\bootstrap\kvcache.go
 * @Description: 注册 go-cachex 缓存（字段级 KV + 整实例）
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package bootstrap

import (
	"github.com/<% .OrgName %>/<% .ProjectName %>/constants"
	cachex "github.com/kamalyes/go-cachex"
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
<%- if eq .GenerateExample "true" %>
	"github.com/<% .OrgName %>/<% .ProjectName %>/models"
<%- end %>
)

// registerKvCache 注册缓存
//
// ============================================================================
// 两个组件，解决两类不同的问题，不要混用
// ============================================================================
//
// ModelKVCache（字段级，NewModelKV）：缓存「主键 → 某个字段」，如 id→name、id→code。
//   专治列表接口的 N+1 回显：拿到一页 100 条记录里的 100 个 tenant_id，
//   要显示租户名，不缓存就是 100 次查询或一次 IN 查询加一轮拼装。
//   读用 GetFieldMany(ctx, "Name", ids) 一次拿回 map[id]name。
//
//   代价：每个缓存字段都注册了一个全量 loader，到 KVCacheRefreshInterval 就
//   SELECT 整张表重建本地映射。所以只给「有界的字典表」注册——租户、仓库、
//   分类这种总量可控的。给会无限增长的业务表注册，等于每隔几分钟扫一次大表。
//
// ModelCache（整实例，NewModelCache）：缓存「复合业务键 → 整条记录」。
//   读走本地 → Redis Hash → 按键精准回源 DB 三层兜底，不做全表加载，
//   所以对大表也安全。适合「按 (tenant_id, id) 取整条配置/实体」这类热点单条读。
//
//   代价：本地那层只有 TTL 淘汰，没有条数上限。键空间极大的表会把内存吃掉，
//   这种场景要么别缓存整实例，要么把 TTL 压短。
//
// 两者都带写后失效广播：Set/Delete 会通过 PubSub 通知其他副本删本地缓存。
// 这是自己用 map 做本地缓存最容易漏的一环——单副本测不出来，
// 多副本上线后表现为「改完有的节点生效有的不生效」。
//
// [EN] Two components for two different problems. ModelKVCache caches primary-key-to-field
// pairs such as id-to-name and exists to kill the N+1 lookup on list endpoints, read in bulk
// with GetFieldMany; its cost is a full-table SELECT per cached field on every refresh tick, so
// it suits bounded reference tables only. ModelCache caches whole rows by a composite business
// key with local, Redis and precise per-key database fallback and no full load, so it is safe
// for large tables; its cost is a local layer bounded only by TTL, not by entry count. Both
// broadcast invalidation over PubSub on write, which is the part a hand-rolled local map
// always misses: it cannot be reproduced on one replica and shows up in production as an update
// taking effect on some nodes but not others.
//
// ============================================================================
// 注册名冲突
// ============================================================================
// 同一张表不能注册两套 ModelKV：注册名只按 表名 + 字段名 生成，不受 namespace 隔离。
// 需要为同一张表建另一套索引（比如用 code 反查而不是用主键），
// 要用 cachex.RegisterKV 自己起一个唯一名字，别再 NewModelKV 一次。
//
// [EN] One table cannot carry two ModelKV registrations: the registration name is derived from
// table plus field and is not isolated by namespace. A second index over the same table, such
// as a lookup by code rather than by primary key, needs cachex.RegisterKV under a name of its
// own.
func (s *<% .ServiceName %>Gateway) registerKvCache() {
	// Redis 没配就直接返回：RegisterKV 在全局客户端为 nil 时会 panic，
	// 而「本地起服务不开 Redis」是很常见的开发姿势，不该因此起不来。
	// 读路径用的是非 panic 的 GetModelKV/GetModelCache，拿不到就退回直查数据库，
	// 所以跳过注册只是失去缓存，功能不受影响。
	//
	// [EN] Return early when Redis is absent: RegisterKV panics on a nil global client, and
	// running locally without Redis is common enough that it must not prevent startup. The read
	// paths use the non-panicking GetModelKV/GetModelCache and fall back to a direct query, so
	// skipping registration costs the cache and nothing else.
	if gwglobal.REDIS == nil {
		gwglobal.LOGGER.Warn("Redis 未配置，跳过缓存注册（将全部直查数据库）| Redis not configured, cache registration skipped")
		return
	}

	cachex.SetGlobalRedisClient(gwglobal.REDIS)
	cachex.SetLogger(gwglobal.LOGGER)
<%- if eq .GenerateExample "true" %>

	// ------------------------------------------------------------------
	// <% .MixedModule %> 模块：聚合根是有界的（租户总量可控），既做字段级又做整实例
	// [EN] The <% .MixedModule %> aggregate root is bounded, so both cache kinds apply
	// ------------------------------------------------------------------
	<% .MixedModule %>Base := cachex.NewModelKVBase().
		Namespace(constants.KVNamespace<% .MixedModule | pascal %>).
		TTL(constants.TTLModelCache).
		RefreshInterval(constants.KVCacheRefreshInterval)

	// id → name / code：跨模块回显与存在性校验都走这里，省掉一次查询
	// 主键由 gorm schema 自动识别，无需 KeyField
	// [EN] id to name and code, used for cross-module display and existence checks; the primary
	// key is taken from the gorm schema
	cachex.NewModelKV[models.<% .MixedService %>Model](<% .MixedModule %>Base, gwglobal.DB).
		Fields("Name", "Code").
		Register()

	// (tenant_id, id) → 整条成员记录
	cachex.NewModelCache[models.<% .MixedService %>MemberModel](<% .MixedModule %>Base, gwglobal.DB).
		KeyFields("TenantId", "Id").
		Register()

	// ------------------------------------------------------------------
	// <% .SecondModule %> 模块：按租户分组的字典表，总量有界
	// [EN] A per-tenant reference table, bounded in size
	// ------------------------------------------------------------------
	<% .SecondModule %>Base := cachex.NewModelKVBase().
		Namespace(constants.KVNamespace<% .SecondModule | pascal %>).
		TTL(constants.TTLModelCache).
		RefreshInterval(constants.KVCacheRefreshInterval)

	cachex.NewModelKV[models.<% .SecondService %>Model](<% .SecondModule %>Base, gwglobal.DB).
		Fields("Name", "Code").
		Register()

	cachex.NewModelCache[models.<% .SecondService %>Model](<% .SecondModule %>Base, gwglobal.DB).
		KeyFields("TenantId", "Id").
		Register()

	// ------------------------------------------------------------------
	// <% .ExampleModule %> 模块：业务主表，会随业务无限增长
	//
	// 这里故意只注册 ModelCache，不注册 ModelKV：
	// ModelKV 的全量刷新会每隔 KVCacheRefreshInterval 扫一次整张表，
	// 对一张持续增长的表来说这是纯粹的浪费，而且随数据量恶化。
	// 需要按 id 批量取名字时，老老实实走一次 IN 查询——
	// 一次有界的查询比定期全表扫描便宜得多。
	//
	// [EN] Deliberately ModelCache only, no ModelKV: the full refresh would scan this growing
	// table on every tick, which is pure waste and gets worse with size. When names are needed
	// in bulk here, one bounded IN query is far cheaper than a periodic full scan.
	// ------------------------------------------------------------------
	<% .ExampleModule %>Base := cachex.NewModelKVBase().
		Namespace(constants.KVNamespace<% .ExampleModule | pascal %>).
		TTL(constants.TTLModelCache)

	cachex.NewModelCache[models.<% .ExampleService %>Model](<% .ExampleModule %>Base, gwglobal.DB).
		KeyFields("TenantId", "Id").
		Register()
<%- else %>

	// TODO: 按需注册本服务的缓存
	//
	// 字段级（有界字典表，列表回显用）：
	//   base := cachex.NewModelKVBase().
	//       Namespace("<% .ProjectName %>:xxx").
	//       TTL(constants.TTLModelCache).
	//       RefreshInterval(constants.KVCacheRefreshInterval)
	//   cachex.NewModelKV[models.XxxModel](base, gwglobal.DB).Fields("Name", "Code").Register()
	//
	// 整实例（按复合键取整条，大表也安全）：
	//   cachex.NewModelCache[models.XxxModel](base, gwglobal.DB).
	//       KeyFields("TenantId", "Id").
	//       Register()
<%- end %>

	gwglobal.LOGGER.Info("Cache registered")
}
