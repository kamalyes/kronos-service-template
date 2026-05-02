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
<%- if eq .GenerateExample "true" %>
	<% .ExampleModule %>svc "github.com/<% .OrgName %>/<% .ProjectName %>/service/<% .ExampleModule %>"
	<% .SecondModule %>svc "github.com/<% .OrgName %>/<% .ProjectName %>/service/<% .SecondModule %>"
	<% .MixedModule %>svc "github.com/<% .OrgName %>/<% .ProjectName %>/service/<% .MixedModule %>"
<%- end %>
)

// <% .ServiceName %>Gateway 服务启动器
type <% .ServiceName %>Gateway struct {
	gateway *gateway.Gateway
	repo    *repository.RepositoryFactory
	svc     *service.<% .ServiceName %>ServiceImpl
<%- if eq .GenerateExample "true" %>
	// exampleSvc 实现协议中心的 <% .ExampleService %>ServiceServer
	exampleSvc *<% .ExampleModule %>svc.<% .ExampleService %>Service
	secondSvc  *<% .SecondModule %>svc.<% .SecondService %>Service
	mixedSvc   *<% .MixedModule %>svc.<% .MixedService %>Service
	mixedMemberSvc *<% .MixedModule %>svc.<% .MixedService %>MemberService
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

	// 先配日志上下文字段，之后所有启动日志和业务日志都带齐链路字段
	// [EN] Configure the log context fields first so every later log carries them
	configureLoggerContext()

	gwglobal.LOGGER.Info("Registering cache...")
	s.registerKvCache()

	gwglobal.LOGGER.Info("Initializing gRPC clients...")
	s.initClients()

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
