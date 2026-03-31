/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\bootstrap\server.go
 * @Description: 服务启动引导，初始化并启动gRPC和HTTP网关服务
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package bootstrap

import (
	"fmt"
	gruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/<% .OrgName %>/<% .ProjectName %>/repository"
	"github.com/<% .OrgName %>/<% .ProjectName %>/service"
	goconfig "github.com/kamalyes/go-config"
	gateway "github.com/kamalyes/go-rpc-gateway"
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
)

// <% .ServiceName %>Gateway 服务启动器
type <% .ServiceName %>Gateway struct {
	gateway *gateway.Gateway
	repo    *repository.RepositoryFactory
	svc     *service.<% .ServiceName %>ServiceImpl
<%- if eq .GenerateExample "true" %>
	// exampleSvc 实现协议中心的 <% .ExampleService %>ServiceServer
	exampleSvc *service.<% .ExampleService %>Service
	secondSvc  *service.<% .SecondService %>Service
	mixedSvc   *service.<% .MixedService %>Service
	mixedMemberSvc *service.<% .MixedService %>MemberService
<%- end %>
}

// New<% .ServiceName %>Gateway 创建服务启动器实例
func New<% .ServiceName %>Gateway() *<% .ServiceName %>Gateway {
	return &<% .ServiceName %>Gateway{}
}

// Run 启动服务
func (s *<% .ServiceName %>Gateway) Run() error {
	gw, err := gateway.NewGateway().
		WithSearchPath("resources").
		WithEnvironment(goconfig.GetEnvironment()).
		WithPrefix("<% .GatewayPrefix %>").
		WithHotReload(nil).
		Build()
	if err != nil {
		return fmt.Errorf("create gateway failed: %v", err)
	}
	s.gateway = gw

	gwglobal.LOGGER.Info("Initializing database...")
	if err := s.initDatabase(); err != nil {
		return err
	}

	gwglobal.LOGGER.Info("Initializing services...")
	s.initServices()

	gwglobal.LOGGER.Info("Registering gRPC services...")
	s.registerGRPCServices()

	gwglobal.LOGGER.Info("Adding gRPC-Gateway middleware...")
	s.gateway.AddGrpcGatewayMiddlewareProvider(s.getMiddleware)

	gwglobal.LOGGER.Info("Rebuilding HTTP gateway...")
	if err := s.gateway.RebuildHTTPGateway(); err != nil {
		return fmt.Errorf("rebuild http gateway failed: %v", err)
	}

	gwglobal.LOGGER.Info("Registering HTTP handlers...")
	s.registerHTTPHandlers()
<%- if eq .UseSubscriber "true" %>

	gwglobal.LOGGER.Info("Starting subscribers...")
	s.startSubscribers()
<%- end %>

	gwglobal.LOGGER.Info("<% .ServiceName %> Server is ready")
	return s.gateway.Run()
}

// getMiddleware 获取中间件列表
func (s *<% .ServiceName %>Gateway) getMiddleware() []gruntime.Middleware {
	return nil
}
