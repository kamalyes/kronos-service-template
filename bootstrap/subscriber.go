/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\bootstrap\subscriber.go
 * @Description: 后台任务与事件订阅启动
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package bootstrap

import (
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
<%- if eq .GenerateExample "true" %>
	<% .MixedModule %>svc "github.com/<% .OrgName %>/<% .ProjectName %>/service/<% .MixedModule %>"
<%- end %>
)

// startSubscribers 启动后台任务与事件订阅
//
// 所有后台 goroutine 都用 s.gateway.Context() 作为父 ctx，不要用
// context.Background()：前者在收到 SIGTERM 时会被取消，goroutine 跟着退出，
// 进程能干净关闭；后者会让 goroutine 一直活着，容器要等到 grace period
// 结束才被 SIGKILL，滚动发布时每个 Pod 都多等几十秒。
//
// [EN] Every background goroutine takes s.gateway.Context() as its parent rather than
// context.Background(): the former is cancelled on SIGTERM so the goroutines exit and the
// process shuts down cleanly, while the latter keeps them alive until the container's grace
// period expires, adding that delay to every pod of a rolling release.
func (s *<% .ServiceName %>Gateway) startSubscribers() {
<%- if eq .GenerateExample "true" %>
	// 下游同步 worker：服务间 gRPC 调用示例
	// 为什么同步放在后台而不是 RPC 里，见 worker 文件开头的说明
	// [EN] The downstream sync worker, an example of a service-to-service gRPC call; the file
	// header explains why the call belongs in the background rather than in the RPC
	<% .MixedModule %>svc.New<% .MixedService %>SyncWorker(s.repo.<% .MixedService %>).
		Start(s.gateway.Context())
<%- else %>
	// TODO: 添加后台任务与事件订阅
	// [EN] Add background tasks and event subscriptions
<%- end %>
	gwglobal.LOGGER.Info("All subscribers started")
}
