/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\models\idgen.go
 * @Description: ID 生成工具函数 | ID generation helpers
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package models

import (
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
	"github.com/kamalyes/go-toolbox/pkg/safe"
)

// 为什么放在 models 而不是 service | Why this lives in models, not service
//
// 分配主键的时机横跨三层：service 创建实体时要、repository 批量插入时要、
// model 的 BeforeCreate 钩子里也可能要。依赖方向是 service → repository → models，
// 所以它只能待在最底层，放在 service 会让下面两层用不了（import 循环）。
//
// [EN] Assigning a primary key happens across all three layers: in service when creating an
// entity, in repository on bulk insert, and possibly in a model's BeforeCreate hook. The
// dependency direction is service → repository → models, so this must sit at the bottom;
// in service the two lower layers could not use it without an import cycle.

// NewSnowflakeID 生成一个新的短 ID 字符串
//
// 委托给网关的全局生成器，不自建：网关启动时已按 osx.GetWorkerIdForSnowflake() 与
// osx.GetDatacenterId() 初始化过一个 Snowflake 生成器（启动日志里的「正在初始化:
// Snowflake」）。再建一个参数相同的生成器，两者各有独立的序列计数器，
// 同一毫秒内会发出相同的 ID。
//
// [EN] Delegates to the gateway's global generator instead of building one. The gateway
// already initializes a Snowflake generator from the same worker and datacenter IDs at
// startup; a second generator with identical parameters keeps its own sequence counter and
// can emit the same ID within the same millisecond.
func NewSnowflakeID() string {
	return gwglobal.NewSnowflakeID()
}

// NewDeterministicID 由业务键生成确定性 ID，相同 key 恒得相同结果
//
// 用于需要幂等的场景：初始化内置数据、按「模块:资源:类型」生成固定节点 ID 等。
// 这类记录每次部署都要落成同一个 ID，不能用 Snowflake（它每次都不同）。
//
// [EN] For idempotent cases such as seeding built-in rows or deriving a stable node ID from
// "module:resource:type". Those records must land on the same ID on every deployment, so
// Snowflake (which differs every call) is not usable there.
func NewDeterministicID(key string) string {
	return safe.ShortHash(key)
}
