<%- if eq .GenerateExample "true" %>
/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\repository\<% .ExampleName | snake %>_repository.go
 * @Description: <% .ExampleName %> 示例仓库，提供数据库操作封装
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package repository

import (
	"context"
	"github.com/<% .OrgName %>/<% .ProjectName %>/models"
	sqlrepo "github.com/kamalyes/go-sqlbuilder/repository"
)

// <% .ExampleName %>Repository 示例仓库结构体，继承通用基础仓库
type <% .ExampleName %>Repository struct {
	*sqlrepo.BaseRepository[models.<% .ExampleName %>Model]
}

// New<% .ExampleName %>Repository 创建示例仓库实例
func New<% .ExampleName %>Repository() *<% .ExampleName %>Repository {
	return &<% .ExampleName %>Repository{
		BaseRepository: newBaseRepository[models.<% .ExampleName %>Model](models.<% .ExampleName %>Model{}.TableName()),
	}
}

// GetByID 根据ID获取记录
func (r *<% .ExampleName %>Repository) GetByID(ctx context.Context, id string) (*models.<% .ExampleName %>Model, error) {
	return r.GetByFilter(ctx, sqlrepo.NewEqFilter("<% .ExampleName | snake %>_id", id))
}

// GetByName 根据名称获取记录
func (r *<% .ExampleName %>Repository) GetByName(ctx context.Context, name string) (*models.<% .ExampleName %>Model, error) {
	return r.GetByFilter(ctx, sqlrepo.NewEqFilter("name", name))
}

// ListWithPagination 分页查询列表
func (r *<% .ExampleName %>Repository) ListWithPagination(ctx context.Context, query *sqlrepo.Query, pagination *sqlrepo.Pagination32) ([]*models.<% .ExampleName %>Model, *sqlrepo.Pagination32, error) {
	return r.ListWithPagination32(ctx, query, pagination)
}

// UpdateByID 根据ID更新指定字段
func (r *<% .ExampleName %>Repository) UpdateByID(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.UpdateFieldsByFilters(ctx, updates, sqlrepo.NewEqFilter("<% .ExampleName | snake %>_id", id))
}

// UpdateStatus 更新状态
func (r *<% .ExampleName %>Repository) UpdateStatus(ctx context.Context, id string, status int32) error {
	return r.UpdateFieldsByFilters(ctx, map[string]interface{}{"status": status}, sqlrepo.NewEqFilter("<% .ExampleName | snake %>_id", id))
}

// NameExists 检查名称是否已存在
func (r *<% .ExampleName %>Repository) NameExists(ctx context.Context, name string) (bool, error) {
	count, err := r.Count(ctx, sqlrepo.NewEqFilter("name", name))
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DeleteByID 根据ID删除记录
func (r *<% .ExampleName %>Repository) DeleteByID(ctx context.Context, id string) error {
	return r.DeleteByFilters(ctx, sqlrepo.NewEqFilter("<% .ExampleName | snake %>_id", id))
}
<%- end %>
