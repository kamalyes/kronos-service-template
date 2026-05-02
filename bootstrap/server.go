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
	"github.com/<% .OrgName %>/<% .ProjectName %>/constants"
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
	// 配置前缀取 constants.ServiceName，不写字面量
	//
	// go-config 按 <前缀>-<环境别名>.yaml 匹配配置文件，所以前缀和文件名必须一致。
	// 两处各写一个字符串的话，改了一处忘了另一处就是启动即崩：
	//   ❌ 未找到前缀为 'xxx' 的配置文件
	// 而这个报错只在真正启动时出现，编译期毫无提示。
	// 都从同一个常量来，这件事结构上就不可能分叉——
	// resources/ 下的文件名也由同一个变量渲染（见 kopy.yaml 的 __project_name__）。
	//
	// [EN] go-config matches <prefix>-<env alias>.yaml, so the prefix and the file name must
	// agree. With a literal in each place, changing one and forgetting the other fails at
	// startup with "config file with prefix not found" and nothing at compile time. Taking both
	// from one constant makes divergence structurally impossible, since the file name under
	// resources/ is rendered from the same variable.
	gw, err := gateway.NewGateway().
		WithSearchPath("resources").
		WithEnvironment(goconfig.GetEnvironment()).
		WithPrefix(constants.ServiceName).
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

	gwglobal.LOGGER.Info("Initializing ID dispenser...")
	s.initIDDispenser()

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
