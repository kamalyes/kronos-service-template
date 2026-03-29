/*
 * @Author: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @Date: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @LastEditors: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @LastEditTime: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @FilePath: \[[ .ProjectName ]]\bootstrap\services.go
 * @Description: 服务初始化和依赖注入
 *
 * Copyright (c) [[ .Year ]] by [[ .AuthorName ]], All Rights Reserved.
 */
package bootstrap

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/[[ .OrgName ]]/[[ .ProjectName ]]/repository"
	"github.com/[[ .OrgName ]]/[[ .ProjectName ]]/service"
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
	"google.golang.org/grpc"
)

// initServices 初始化所有服务
func (s *[[ .ServiceName ]]Gateway) initServices() {
	gwglobal.LOGGER.Info("Initializing services...")

	s.repo = repository.NewRepositoryFactory()
	s.svc = service.New[[ .ServiceName ]]Service(s.repo)

	gwglobal.LOGGER.Info("Services initialized")
}

// registerGRPCServices 注册gRPC服务
func (s *[[ .ServiceName ]]Gateway) registerGRPCServices() {
	s.gateway.RegisterService(func(srv *grpc.Server) {
		// TODO: 注册 gRPC 服务
		// pb.Register[[ .ServiceName ]]ServiceServer(srv, s.svc)
	})
	gwglobal.LOGGER.Info("[[ .ServiceName ]]Service gRPC registered")
}

// registerHTTPHandlers 注册HTTP处理器
func (s *[[ .ServiceName ]]Gateway) registerHTTPHandlers() {
	// TODO: 注册 HTTP 处理器
	// if err := s.gateway.RegisterGatewayHandler(func(ctx context.Context, mux *runtime.ServeMux) error {
	// 	return pb.Register[[ .ServiceName ]]ServiceHandlerServer(ctx, mux, s.svc)
	// }); err != nil {
	// 	gwglobal.LOGGER.Fatal("Failed to register [[ .ServiceName ]]Service HTTP handler: %v", err)
	// }

	_ = context.Background()
	_ = &runtime.ServeMux{}

	gwglobal.LOGGER.Info("[[ .ServiceName ]]Service HTTP handler registered")
}
