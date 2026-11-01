/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 17:05:08
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 19:28:30
 * @FilePath: \<% .ProjectName %>\constants\service_name.go
 * @Description: 服务名称常量
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */

package constants

// ServiceName 服务名称
const ServiceName = "<% .ProjectName %>"

// 远程服务名 | Remote service names
//
// 这里的字符串必须与 resources/<% .ProjectName %>-dev.yaml 里 grpc.clients
// 的 key 完全一致——连接池是按服务名索引的，对不上就是「配了也拿不到连接」，
// 而且 InitClient 只返回 ok=false，没有任何报错，表现为那条调用静默退化成
// 「下游未配置」。所以两边都引用同一个常量，别在任何一侧写字面量。
//
// [EN] These strings must match the keys under grpc.clients in the resources file exactly: the
// connection pool is indexed by service name, and a mismatch means the client is configured yet
// never resolves. InitClient merely returns false with no error, so the call silently degrades
// to "downstream not configured". Reference this constant on both sides rather than writing the
// literal anywhere.
const (
	// RemoteServiceKey 中心发号 / 发密钥服务
	//
	// 由 kronos-key-service-template 生成，实现协议中心 key 模块的
	// IdBlockService（ID 号段）与 TempKeyService（临时 RSA 密钥）。
	// 本服务目前只用到前者，见 bootstrap/idgen.go。
	//
	// [EN] Generated from kronos-key-service-template and implementing the contract hub's key
	// module: IdBlockService for id blocks and TempKeyService for temporary RSA keys. Only the
	// former is used here, in bootstrap/idgen.go.
	RemoteServiceKey = "<% .DemoRemoteService %>"

	// SecondRemoteServiceKey 第二下游，实现协议中心 <% .SecondModule %> 模块的
	// <% .SecondService %>Service
	//
	// 由混合模块的同步 worker 调用，用来演示服务间 gRPC 调用，
	// 见 service/<% .MixedModule %>/<% .MixedService | snake %>_sync_worker.go
	//
	// [EN] The second downstream implementing the contract hub's <% .SecondService %>Service in
	// the <% .SecondModule %> module, called by the mixed module's sync worker as a
	// service-to-service example
	SecondRemoteServiceKey = "<% .DemoSecondRemoteService %>"
)
