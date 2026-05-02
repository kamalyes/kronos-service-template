/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\bootstrap\idgen.go
 * @Description: ID 中心分发器装配
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package bootstrap

import (
	"context"

	"github.com/<% .OrgName %>/<% .ProjectName %>/constants"
	"github.com/<% .OrgName %>/<% .ProjectName %>/models"
	cachex "github.com/kamalyes/go-cachex"
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
)

// initIDDispenser 装配 ID 中心分发器
//
// 为什么要有这一步，见 models/idsource.go 开头那段：本地 Snowflake 的 worker id
// 在 Deployment 下是哈希取模出来的，而槽位只有 32 个，8 个副本就有 61% 的概率
// 有两个 Pod 撞上同一个 worker id，撞上之后同一毫秒内会发出相同的 ID。
//
// 这里用 go-cachex 的 Redis IDGenerator 当分发器：它的 GenerateIDs 是一条
// Lua 脚本里 INCRBY n，一次往返拿一整段，原子性由 Redis 单线程保证，
// 不需要额外的锁。Redis 是所有实例共享的，于是「谁发号」这件事收敛到一处，
// 实例数、Pod 名、worker id 全都不再参与 ID 的唯一性。
//
// Redis 没配就不装配：models.NewID 此时等价于 NewSnowflakeID，
// 模板开箱即用不需要 Redis。单实例开发环境下 Snowflake 本来就不会撞。
//
// 要换成独立的发号服务（比如一个专门的 id-dispenser gRPC 服务），
// 改的只有下面那个闭包——把 GenerateIDs 换成一次 gRPC 调用即可，
// 客户端怎么接见 service/common/client.go。models 侧一行都不用动。
//
// [EN] The reason for this step is at the top of models/idsource.go: a Deployment derives the
// Snowflake worker id from a hostname hash into only 32 slots, so eight replicas carry a 61%
// chance that two pods share one and then emit identical ids within a millisecond. The dispenser
// here is go-cachex's Redis IDGenerator, whose GenerateIDs performs INCRBY n inside a Lua script
// and so returns a whole block in one atomic round trip, with Redis's single-threaded execution
// standing in for a lock. Redis is shared by every instance, which collapses the question of who
// issues ids to one place and removes instance count, pod name and worker id from the
// uniqueness argument entirely. Without Redis nothing is assembled and NewID behaves as
// NewSnowflakeID, so the template runs out of the box; a single-instance dev environment cannot
// collide anyway. Swapping in a dedicated dispenser service means replacing the closure below
// with a gRPC call and changing nothing under models.
func (s *<% .ServiceName %>Gateway) initIDDispenser() {
	if gwglobal.REDIS == nil {
		gwglobal.LOGGER.Warn("Redis 未配置，ID 使用本地 Snowflake" +
			" | Redis not configured, ids come from local snowflake")
		return
	}

	// keyPrefix 必须带服务名：多个服务共用一套 Redis 时，前缀相同就等于
	// 共用同一个序列，两个服务各自的号段会交错——不会重复（序列是同一个），
	// 但一个服务的写入量会推高另一个服务的号，而且谁都查不出为什么。
	//
	// timeFormat 决定序列多久归零一次。默认按小时（2006010215），
	// 配合 expire 让旧的序列键自动清掉，不会无限堆积。
	//
	// [EN] The key prefix must carry the service name: with several services on one Redis an
	// identical prefix means one shared sequence, so their blocks interleave. They do not
	// collide, since it is the same sequence, but one service's write volume pushes the other's
	// numbers along and nobody can tell why. The time format sets how often the sequence resets,
	// hourly by default, and the expiry clears old sequence keys instead of letting them pile up.
	gen := cachex.NewIDGenerator(gwglobal.REDIS).
		WithKeyPrefix(constants.ServiceName + ":idgen:").
		WithSequenceLength(constants.IDSequenceLength).
		WithExpire(constants.IDSequenceExpire)

	models.SetIDAllocator(func(ctx context.Context, n int) ([]string, error) {
		return gen.GenerateIDs(ctx, n)
	}, constants.IDBlockSize)

	gwglobal.LOGGER.Info("ID dispenser initialized (block size %d)", constants.IDBlockSize)
}
