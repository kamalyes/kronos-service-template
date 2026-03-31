/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\repository\base_cached.go
 * @Description: 带缓存与事务的泛型仓储基座，所有实体仓库共用
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package repository

import (
	"context"
	"encoding/json"
	"time"

	cachex "github.com/kamalyes/go-cachex"
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
	sqlconstants "github.com/kamalyes/go-sqlbuilder/constants"
	sqlrepo "github.com/kamalyes/go-sqlbuilder/repository"
)

// DefaultCacheTTL 单条详情的缓存时长
//
// 比配置里的 cache.default-ttl（24h）短得多：详情是会被改的，TTL 是失效漏掉时的兜底，
// 不是主要失效手段。主要手段是写路径显式删键。
// TTL 定得长会让一次漏删的影响持续很久，定得短则缓存形同虚设，分钟级是常见折中。
//
// [EN] Much shorter than the configured cache.default-ttl: a detail row changes, so the TTL is
// the backstop for a missed invalidation, not the primary mechanism — explicit deletion on the
// write path is. A long TTL makes one missed deletion linger; a very short one makes the cache
// pointless.
const DefaultCacheTTL = 5 * time.Minute

// TenantScoped 要求实体能报出自己的主键与所属租户
//
// 泛型基座靠它定位缓存键与校验归属，所有落库实体都应实现这两个方法
// [EN] The generic base uses this to build the cache key and verify ownership
type TenantScoped interface {
	GetID() string
	GetTenantID() string
}

// CachedRepository 在 go-sqlbuilder 的 BaseRepository 之上加缓存与事务批量操作
//
// 为什么做成泛型基座而不是每个实体各写一份：缓存键怎么拼、什么时候失效、
// 批量删除要不要事务，这些决定一旦在四个仓库里各写一遍，就会在第五个实体上走样。
// 收在一处后，新增实体只需声明 entity 名与表名。
//
// [EN] A generic base rather than one copy per entity: how the cache key is built, when it is
// invalidated and whether batch deletes are transactional are decisions that drift once copied
// into four repositories. Centralized, a new entity only declares its name and table.
type CachedRepository[T any] struct {
	*sqlrepo.BaseRepository[T]
	// entity 缓存键的命名空间，取实体名的蛇形形式
	entity string
}

// NewCachedRepository 构造泛型仓储基座
func NewCachedRepository[T any](entity, tableName string) *CachedRepository[T] {
	return &CachedRepository[T]{
		BaseRepository: newBaseRepository[T](tableName),
		entity:         entity,
	}
}

// cacheKey 构造缓存键
//
// 键里必须带 tenantID：缓存键是另一个越权读的入口，只用 id 做键的话，
// A 租户读过的记录会被 B 租户命中。配置里的 cache.key-prefix 由缓存实现自己加，
// 这里不重复拼。
//
// [EN] The key must include the tenant: a cache key is another cross-tenant read hole, since an
// id-only key lets tenant B hit what tenant A loaded. The configured key-prefix is added by the
// cache implementation itself and is not repeated here.
func (r *CachedRepository[T]) cacheKey(tenantID, id string) []byte {
	return []byte(r.entity + ":" + tenantID + ":" + id)
}

// cache 取网关构建好的缓存实例
//
// 缓存实例由网关的 PoolManager 按 resources 里的 cache 段构建（type / enabled /
// default-ttl / key-prefix 都在那里），服务不自己建连接。
// 配置里 enabled 为 false 或连接池未初始化时返回 nil，调用方自动退化为直查数据库。
//
// [EN] The cache is built by the gateway's PoolManager from the cache section in resources, so
// the service never creates its own connection. It returns nil when disabled or uninitialized,
// and callers degrade to a direct query.
func (r *CachedRepository[T]) cache() *cachex.CtxCache {
	if gwglobal.POOL_MANAGER == nil {
		return nil
	}
	c := gwglobal.POOL_MANAGER.GetCache()
	return &c
}

// GetByIDScoped 按主键查询并校验租户归属
//
// 所有单条查询都带 tenantID 过滤：多租户服务里只按主键查是越权读的经典口子，
// 调用方只要猜到别的租户的 ID 就能读到数据。tenantID 由拦截器从请求上下文解析，
// 不从请求体取（请求体是调用方可控的）。
//
// [EN] Every single-row read is tenant-scoped: looking up by primary key alone is the classic
// cross-tenant read hole. The tenant comes from the request context, never from the request
// body, which the caller controls.
func (r *CachedRepository[T]) GetByIDScoped(ctx context.Context, tenantID, id string) (*T, error) {
	return r.GetByFilters(ctx,
		sqlrepo.NewEqFilter("id", id),
		sqlrepo.NewEqFilter("tenant_id", tenantID),
	)
}

// GetByIDCached 带缓存的单条查询
//
// 用 GetOrCompute 而不是「先 Get 再 Set」：后者在缓存未命中且并发量大时，
// 所有请求会同时穿到数据库（缓存击穿）。GetOrCompute 内部用 singleflight 去重，
// 同一个键并发未命中时只有一个请求执行 loader，其余等它的结果。
//
// 未找到的记录不写缓存：写了就要再处理空值穿透与空值失效两件事，
// 而「不存在」通常不是高频查询，收益不抵复杂度。
//
// [EN] GetOrCompute rather than get-then-set: the latter lets every concurrent miss hit the
// database at once. GetOrCompute deduplicates with singleflight, so one request runs the loader
// and the rest wait on it. A missing row is not cached: that would add negative caching and its
// own invalidation, and "not found" is rarely the hot path.
func (r *CachedRepository[T]) GetByIDCached(ctx context.Context, tenantID, id string) (*T, error) {
	c := r.cache()
	if c == nil {
		return r.GetByIDScoped(ctx, tenantID, id)
	}

	raw, err := c.GetOrCompute(ctx, r.cacheKey(tenantID, id), DefaultCacheTTL,
		func(ctx context.Context) ([]byte, error) {
			entity, err := r.GetByIDScoped(ctx, tenantID, id)
			if err != nil {
				return nil, err
			}
			if entity == nil {
				// 返回 nil 让 GetOrCompute 跳过写缓存
				// [EN] Returning nil makes GetOrCompute skip the cache write
				return nil, nil
			}
			return json.Marshal(entity)
		})
	if err != nil || len(raw) == 0 {
		// 缓存层出错不该让查询失败：退回直查数据库
		// [EN] A cache failure must not fail the query; fall back to a direct read
		return r.GetByIDScoped(ctx, tenantID, id)
	}

	var entity T
	if err := json.Unmarshal(raw, &entity); err != nil {
		return r.GetByIDScoped(ctx, tenantID, id)
	}
	return &entity, nil
}

// Invalidate 删除单条缓存
//
// 写路径一律「改完就删键」，不是「改完就写新值」：写新值要求写缓存与写库原子，
// 否则并发更新下两者会错位（A 写库 → B 写库 → B 写缓存 → A 写缓存，缓存留下旧值）。
// 删键的失败模式只是下一次读多查一次库。
//
// [EN] The write path deletes the key rather than writing the new value: writing through would
// need cache and database updated atomically, or concurrent updates leave a stale value
// (A writes db, B writes db, B writes cache, A writes cache). A failed delete merely costs one
// extra query on the next read.
func (r *CachedRepository[T]) Invalidate(ctx context.Context, tenantID, id string) {
	c := r.cache()
	if c == nil {
		return
	}
	_ = c.Del(ctx, r.cacheKey(tenantID, id))
}

// UpdateScoped 按主键更新并失效缓存
func (r *CachedRepository[T]) UpdateScoped(
	ctx context.Context, tenantID, id string, updates map[string]interface{},
) error {
	if len(updates) == 0 {
		return nil
	}
	if err := r.UpdateFieldsByFilters(ctx, updates,
		sqlrepo.NewEqFilter("id", id),
		sqlrepo.NewEqFilter("tenant_id", tenantID),
	); err != nil {
		return err
	}
	r.Invalidate(ctx, tenantID, id)
	return nil
}

// DeleteScoped 按主键删除并失效缓存
func (r *CachedRepository[T]) DeleteScoped(ctx context.Context, tenantID, id string) error {
	if err := r.DeleteByFilters(ctx,
		sqlrepo.NewEqFilter("id", id),
		sqlrepo.NewEqFilter("tenant_id", tenantID),
	); err != nil {
		return err
	}
	r.Invalidate(ctx, tenantID, id)
	return nil
}

// ExistsByField 判断某字段在租户内是否已存在，excludeID 非空时排除自身（供更新时校验）
func (r *CachedRepository[T]) ExistsByField(
	ctx context.Context, tenantID, field, value, excludeID string,
) (bool, error) {
	filters := []*sqlrepo.Filter{
		sqlrepo.NewEqFilter("tenant_id", tenantID),
		sqlrepo.NewEqFilter(field, value),
	}
	if excludeID != "" {
		filters = append(filters, sqlrepo.NewFilter("id", sqlconstants.OP_NEQ, excludeID))
	}
	count, err := r.Count(ctx, filters...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// BatchDeleteByIDs 事务内批量删除，返回实际删除的 ID 列表
//
// 为什么要事务：逐条删除是 N 次往返，且中途失败会留下删一半的状态——
// 调用方拿到「部分成功」后无法判断该重试哪些，重试又可能重复删。
// 放进一个事务后要么全删要么全不删，调用方的重试语义才是确定的。
//
// 先查出存在的记录再删：批量接口的调用方常常混进不属于本租户的 ID，
// 直接按 ID 删会让「删了 0 行」和「删成功」无法区分，也顺手拿到了要失效的键。
//
// [EN] Why a transaction: deleting one by one is N round trips and a mid-way failure leaves a
// half-deleted state, after which the caller cannot tell what to retry. Inside one transaction
// it is all or nothing, so retry semantics are well defined. Existing rows are fetched first
// because batch callers routinely include ids from other tenants, which also yields the keys to
// invalidate.
func (r *CachedRepository[T]) BatchDeleteByIDs(
	ctx context.Context, tenantID string, ids []string,
) ([]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	q := sqlrepo.NewQuery().
		AddFilter(sqlrepo.NewEqFilter("tenant_id", tenantID)).
		AddFilter(sqlrepo.NewFilter("id", sqlconstants.OP_IN, ids))
	existing, err := r.List(ctx, q)
	if err != nil {
		return nil, err
	}
	if len(existing) == 0 {
		return nil, nil
	}

	if err := r.Transaction(ctx, func(tx sqlrepo.Transaction[T]) error {
		return tx.DeleteBatch(ctx, existing...)
	}); err != nil {
		return nil, err
	}

	deleted := make([]string, 0, len(existing))
	for _, item := range existing {
		if scoped, ok := any(item).(TenantScoped); ok {
			deleted = append(deleted, scoped.GetID())
			r.Invalidate(ctx, tenantID, scoped.GetID())
		}
	}
	return deleted, nil
}
