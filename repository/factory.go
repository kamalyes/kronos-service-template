/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\repository\factory.go
 * @Description: 仓库工厂，统一管理所有仓库实例
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package repository

import (
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
	"github.com/kamalyes/go-sqlbuilder/db"
	sqlrepo "github.com/kamalyes/go-sqlbuilder/repository"
	"sync"
)

// RepositoryFactory 仓库工厂结构体
type RepositoryFactory struct {
<%- if eq .GenerateExample "true" %>
	<% .ExampleService %> *<% .ExampleService %>Repository
	<% .SecondService %>  *<% .SecondService %>Repository
	<% .MixedService %>   *<% .MixedService %>Repository
	<% .MixedService %>Member *<% .MixedService %>MemberRepository
<%- else %>
	// TODO: 添加仓库字段
<%- end %>
}

var (
	factoryInstance *RepositoryFactory
	factoryOnce     sync.Once
)

// NewRepositoryFactory 创建仓库工厂单例
func NewRepositoryFactory() *RepositoryFactory {
	factoryOnce.Do(func() {
		factoryInstance = &RepositoryFactory{
<%- if eq .GenerateExample "true" %>
			<% .ExampleService %>: New<% .ExampleService %>Repository(),
			<% .SecondService %>:  New<% .SecondService %>Repository(),
			<% .MixedService %>:   New<% .MixedService %>Repository(),
			<% .MixedService %>Member: New<% .MixedService %>MemberRepository(),
<%- else %>
			// TODO: 初始化仓库
<%- end %>
		}
	})
	return factoryInstance
}

// newBaseRepository 构造通用 BaseRepository 实例
func newBaseRepository[T any](tableName string) *sqlrepo.BaseRepository[T] {
	dbHandler := db.MustNewGormHandler(gwglobal.DB)
	return sqlrepo.NewBaseRepository[T](
		dbHandler,
		gwglobal.GetLogger(),
		tableName,
		sqlrepo.WithAutoFields[T](),
	)
}
