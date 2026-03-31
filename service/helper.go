/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\service\helper.go
 * @Description: 服务层公共辅助：租户解析与日志字段
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package service

import (
	"context"

	bizerror "github.com/kamalyes/go-bizerror"
	gwmiddleware "github.com/kamalyes/go-rpc-gateway/middleware"
)

// tenantOf 从请求上下文取租户 ID
//
// 租户只认上下文，不认请求体：请求体里的 tenant_id 是调用方可控的，拿它过滤等于没过滤。
// 上下文里的值由网关的 request_context 中间件按 resources 配置的
// authorization-sources 解析出来（见 resources/<% .GatewayPrefix %>-dev.yaml）。
//
// [EN] The tenant comes from the context only, never the request body, which the caller
// controls. The gateway's request_context middleware resolves it per the configured sources.
func tenantOf(ctx context.Context) (string, error) {
	tenantID := gwmiddleware.GetTenantID(ctx)
	if tenantID == "" {
		return "", bizerror.ToGRPCError(ctx, bizerror.BizErrCodeUnauthenticated)
	}
	return tenantID, nil
}

// logKV 每条业务日志都带上的公共字段
//
// 搭配 *ContextKV 使用，不要用 *KV：go-logger 的 Context 变体会自动带出 trace_id
// （优先取 OTel span——网关的 gRPC server 挂了 otelgrpc 的 StatsHandler，
// 所以每个 RPC 的 ctx 都在 span 内），用非 Context 变体会把这条链路丢掉。
// 而 request_id / user_id / tenant_id 不在 go-logger 的默认抽取列表里
// （默认只有 trace_id），所以这里用 gwmiddleware 的 getter 显式补上。
//
// [EN] Use with the *ContextKV variants, never *KV: go-logger's Context family emits trace_id
// automatically, preferring the OTel span (the gateway's gRPC server installs otelgrpc's
// StatsHandler, so every RPC context is inside a span). The non-context variants lose it.
// request_id / user_id / tenant_id are not in go-logger's default extraction list — only
// trace_id is — so they are added explicitly through gwmiddleware getters.
func logKV(ctx context.Context, kv ...any) []any {
	base := []any{
		"request_id", gwmiddleware.GetRequestID(ctx),
		"tenant_id", gwmiddleware.GetTenantID(ctx),
		"user_id", gwmiddleware.GetUserID(ctx),
	}
	return append(base, kv...)
}

// dbError 把仓储层的错误统一归并成业务错误码
//
// 不把 err 原样返回给调用方：底层错误里可能带表名、SQL 片段、连接串，
// 那是内部信息。日志里记全文，响应里只给错误码与本地化文案。
//
// [EN] The underlying error is not returned to the caller: it can carry table names, SQL
// fragments or connection strings. The full text goes to the log; the response carries only
// the code and its localized message.
func dbError(ctx context.Context) error {
	return bizerror.ToGRPCError(ctx, bizerror.BizErrCodeDatabaseOperationFailed)
}

// notFound 资源不存在
func notFound(ctx context.Context) error {
	return bizerror.ToGRPCError(ctx, bizerror.BizErrCodeNotFound)
}
