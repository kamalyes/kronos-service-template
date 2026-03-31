/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\service\<% .ExampleName | snake %>_service.go
 * @Description: <% .ServiceName %> 服务实现
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package service

import (
<%- if eq .GenerateExample "true" %>
	"context"

	"github.com/<% .OrgName %>/<% .ProjectName %>/models"
	"github.com/<% .OrgName %>/<% .ProjectName %>/repository"
	sqlrepo "github.com/kamalyes/go-sqlbuilder/repository"
<%- end %>
)

<%- if eq .GenerateExample "true" %>

// <% .ExampleName %>Service <% .ExampleName %> 示例服务
type <% .ExampleName %>Service struct {
	repo *repository.<% .ExampleName %>Repository
}

// New<% .ExampleName %>Service 创建 <% .ExampleName %> 服务实例
func New<% .ExampleName %>Service(repo *repository.<% .ExampleName %>Repository) *<% .ExampleName %>Service {
	return &<% .ExampleName %>Service{repo: repo}
}

// GetByID 根据ID获取<% .ExampleName %>
func (s *<% .ExampleName %>Service) GetByID(ctx context.Context, id string) (*models.<% .ExampleName %>Model, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByName 根据名称获取<% .ExampleName %>
func (s *<% .ExampleName %>Service) GetByName(ctx context.Context, name string) (*models.<% .ExampleName %>Model, error) {
	return s.repo.GetByName(ctx, name)
}

// List 分页查询<% .ExampleName %>列表
func (s *<% .ExampleName %>Service) List(ctx context.Context, query *sqlrepo.Query, pagination *sqlrepo.Pagination32) ([]*models.<% .ExampleName %>Model, *sqlrepo.Pagination32, error) {
	return s.repo.ListWithPagination(ctx, query, pagination)
}

// Create 创建<% .ExampleName %>
func (s *<% .ExampleName %>Service) Create(ctx context.Context, model *models.<% .ExampleName %>Model) error {
	_, err := s.repo.Create(ctx, model)
	return err
}

// Update 更新<% .ExampleName %>
func (s *<% .ExampleName %>Service) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return s.repo.UpdateByID(ctx, id, updates)
}

// UpdateStatus 更新<% .ExampleName %>状态
func (s *<% .ExampleName %>Service) UpdateStatus(ctx context.Context, id string, status int32) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

// Delete 删除<% .ExampleName %>
func (s *<% .ExampleName %>Service) Delete(ctx context.Context, id string) error {
	return s.repo.DeleteByID(ctx, id)
}

// NameExists 检查名称是否已存在
func (s *<% .ExampleName %>Service) NameExists(ctx context.Context, name string) (bool, error) {
	return s.repo.NameExists(ctx, name)
}
<%- end %>
