/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 17:05:08
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 19:28:30
 * @FilePath: \<% .ProjectName %>\bootstrap\services.go
 * @Description: 服务初始化和依赖注入
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package bootstrap

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/<% .OrgName %>/<% .ProjectName %>/repository"
	"github.com/<% .OrgName %>/<% .ProjectName %>/service"
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
	"google.golang.org/grpc"
<%- if eq .GenerateExample "true" %>
	<% .ExampleModule %>svc "github.com/<% .OrgName %>/<% .ProjectName %>/service/<% .ExampleModule %>"
	<% .SecondModule %>svc "github.com/<% .OrgName %>/<% .ProjectName %>/service/<% .SecondModule %>"
	<% .MixedModule %>svc "github.com/<% .OrgName %>/<% .ProjectName %>/service/<% .MixedModule %>"
	pb "<% .ShareProtoModule %>/pb/<% .ExampleModule %>"
	secondpb "<% .ShareProtoModule %>/pb/<% .SecondModule %>"
	mixedpb "<% .ShareProtoModule %>/pb/<% .MixedModule %>"
<%- end %>
)

// initServices 初始化所有服务
func (s *<% .ServiceName %>Gateway) initServices() {
	gwglobal.LOGGER.Info("Initializing services...")

	// 转换器先注册：registerConverters 会按反射构建字段计划并缓存，
	// 放到开始服务之后注册，第一个请求就要替所有人付这笔反射开销
	// [EN] Register converters first; doing it after serving starts makes the first request
	// pay the reflection cost for everyone
	registerConverters()

	s.repo = repository.NewRepositoryFactory()
	s.svc = service.New<% .ServiceName %>Service(s.repo)
<%- if eq .GenerateExample "true" %>
	s.exampleSvc = <% .ExampleModule %>svc.New<% .ExampleService %>Service(s.repo.<% .ExampleService %>)
	s.secondSvc = <% .SecondModule %>svc.New<% .SecondService %>Service(s.repo.<% .SecondService %>)
	s.mixedSvc = <% .MixedModule %>svc.New<% .MixedService %>Service(s.repo.<% .MixedService %>)
	s.mixedMemberSvc = <% .MixedModule %>svc.New<% .MixedService %>MemberService(s.repo.<% .MixedService %>Member)
<%- end %>

	gwglobal.LOGGER.Info("Services initialized")
}

// registerGRPCServices 注册gRPC服务
func (s *<% .ServiceName %>Gateway) registerGRPCServices() {
	s.gateway.RegisterService(func(srv *grpc.Server) {
<%- if eq .GenerateExample "true" %>
		pb.Register<% .ExampleService %>ServiceServer(srv, s.exampleSvc)
		secondpb.Register<% .SecondService %>ServiceServer(srv, s.secondSvc)
		mixedpb.Register<% .MixedService %>ServiceServer(srv, s.mixedSvc)
		mixedpb.Register<% .MixedService %>MemberServiceServer(srv, s.mixedMemberSvc)
<%- else %>
		// TODO: 注册 gRPC 服务
		// pb.RegisterXxxServiceServer(srv, s.xxxSvc)
<%- end %>
	})
	gwglobal.LOGGER.Info("gRPC services registered")
}

// registerHTTPHandlers 注册HTTP处理器
func (s *<% .ServiceName %>Gateway) registerHTTPHandlers() {
<%- if eq .GenerateExample "true" %>
	// 用 HandlerServer 而不是 HandlerFromEndpoint：前者把 HTTP 请求直接交给同进程的
	// 实现，不绕回 gRPC 端口，省掉一次本机网络往返与一次序列化
	// [EN] HandlerServer dispatches HTTP straight to the in-process implementation instead of
	// dialing back into the gRPC port, saving a loopback round trip and one serialization
	if err := s.gateway.RegisterGatewayHandler(func(ctx context.Context, mux *runtime.ServeMux) error {
		if err := pb.Register<% .ExampleService %>ServiceHandlerServer(ctx, mux, s.exampleSvc); err != nil {
			return err
		}
		if err := secondpb.Register<% .SecondService %>ServiceHandlerServer(ctx, mux, s.secondSvc); err != nil {
			return err
		}
		if err := mixedpb.Register<% .MixedService %>ServiceHandlerServer(ctx, mux, s.mixedSvc); err != nil {
			return err
		}
		return mixedpb.Register<% .MixedService %>MemberServiceHandlerServer(ctx, mux, s.mixedMemberSvc)
	}); err != nil {
		gwglobal.LOGGER.Fatal("注册 HTTP handler 失败 | failed to register HTTP handler: %v", err)
	}
<%- else %>
	// TODO: 注册 HTTP 处理器
	// if err := s.gateway.RegisterGatewayHandler(func(ctx context.Context, mux *runtime.ServeMux) error {
	// 	return pb.RegisterXxxServiceHandlerServer(ctx, mux, s.xxxSvc)
	// }); err != nil {
	// 	gwglobal.LOGGER.Fatal("failed to register HTTP handler: %v", err)
	// }
	_ = context.Background()
	_ = &runtime.ServeMux{}
<%- end %>

	gwglobal.LOGGER.Info("HTTP handlers registered")
}
