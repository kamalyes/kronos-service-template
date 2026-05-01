/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\constants\cache.go
 * @Description: 缓存命名空间与时长常量
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */

package constants

import "time"

// 缓存命名空间 | Cache namespaces
//
// 命名空间会拼进 Redis 键，用途是在同一个 Redis 实例上隔离键空间。
// 带上服务名是必须的：多个服务共用一套 Redis 时，只写 "tenant" 这种模块名
// 会让两个服务的同名表互相覆盖，而且覆盖的现象是「读到别的服务的数据」，
// 不是报错，排查起来极其费时。
//
// [EN] The namespace is part of the Redis key and isolates key space on a shared Redis
// instance. Including the service name is not optional: with several services on one Redis, a
// bare module name lets two services with a same-named table overwrite each other, and the
// symptom is reading another service's data rather than an error.
const (
	KVNamespace<% .ExampleModule | pascal %> = "<% .ProjectName %>:<% .ExampleModule %>"
	KVNamespace<% .SecondModule | pascal %> = "<% .ProjectName %>:<% .SecondModule %>"
	KVNamespace<% .MixedModule | pascal %> = "<% .ProjectName %>:<% .MixedModule %>"
)

// 缓存时长 | Cache durations
const (
	// TTLModelCache 整实例缓存的 TTL
	//
	// 写路径会主动失效，TTL 只是漏删时的兜底，所以不必很短；但也不能很长，
	// 否则一次漏删的影响会持续很久。半小时量级是常见取值。
	//
	// [EN] The write path invalidates explicitly, so the TTL is only a backstop for a missed
	// deletion. It need not be short, but a long one makes one missed deletion linger.
	TTLModelCache = 30 * time.Minute

	// KVCacheRefreshInterval 字段级 KV 缓存的全量刷新间隔
	//
	// 注意这个刷新是「全表 SELECT」：go-cachex 的 ModelKVCache 为每个缓存字段
	// 注册一个全量 loader，到点就把整张表扫一遍重建本地映射。
	// 所以只给有界的字典表注册 ModelKV（租户、仓库、分类这种几千行以内的），
	// 不要给会无限增长的业务表注册，否则每个间隔都在扫大表。
	//
	// [EN] This refresh is a full table SELECT: ModelKVCache registers a full loader per cached
	// field and rebuilds the local map on each tick. Register ModelKV only for bounded reference
	// tables, never for a business table that grows without bound.
	KVCacheRefreshInterval = 5 * time.Minute

	// ModelKVLoadTTL 同一次刷新内多个字段共用一次全表查询的去重窗口
	// [EN] The window in which several fields of one refresh share a single full-table query
	ModelKVLoadTTL = 30 * time.Second
)
