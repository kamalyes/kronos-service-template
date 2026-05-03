/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\bootstrap\clients.go
 * @Description: 服务间 gRPC 客户端装配
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package bootstrap

import (
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
<%- if eq .GenerateExample "true" %>
	"github.com/<% .OrgName %>/<% .ProjectName %>/service/common"
<%- end %>
)

// initClients 装配下游 gRPC 客户端
//
// 放在 initServices 之前：service 层通过 common.XxxRemote() 取客户端，
// 取的是装配好的单例。顺序反了的话，service 构造时拿到的是空壳——
// 不过因为取用是在调用时而不是构造时，顺序错了也不会立刻崩，
// 而是表现为「启动后一段时间内所有远程调用都认为下游未配置」，更难查。
//
// 客户端装配失败不阻断启动：下游不可用不该让本服务起不来，
// 否则一个边缘下游的故障会沿调用链反向放大成一片服务起不来。
//
// [EN] Before initServices, because the service layer fetches the assembled singleton through
// common.XxxRemote(). The wrong order does not crash, since the fetch happens at call time
// rather than construction time; it shows up as every remote call believing the downstream is
// unconfigured, which is harder to diagnose. A failure here never blocks startup: an
// unavailable downstream must not keep this service down, or one peripheral failure amplifies
// backwards along the call chain into a set of services that cannot start.
func (s *<% .ServiceName %>Gateway) initClients() {
<%- if eq .GenerateExample "true" %>
	common.SetClients(common.NewClients(s.gateway))
<%- else %>
	// TODO: 有下游依赖时在 service/common/client.go 里声明并在此装配
	// [EN] Declare downstream clients in service/common/client.go and assemble them here
<%- end %>
	gwglobal.LOGGER.Info("gRPC clients initialized")
}
