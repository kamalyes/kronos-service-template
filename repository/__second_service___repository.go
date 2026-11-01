/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 17:05:08
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 19:28:30
 * @FilePath: \<% .ProjectName %>\repository\<% .SecondService | snake %>_repository.go
 * @Description: <% .SecondService %> 数据访问层
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package repository
<%- if eq .GenerateExample "true" %>

import (
	"context"

	"github.com/<% .OrgName %>/<% .ProjectName %>/models"
	pb "<% .ShareProtoModule %>/pb/<% .SecondModule %>"
	sqlconstants "github.com/kamalyes/go-sqlbuilder/constants"
	sqlrepo "github.com/kamalyes/go-sqlbuilder/repository"
)

// <% .SecondService %>SortableFields 允许排序的字段白名单
//
// 这个实体带 sort_order，是业务上真正的默认排序字段，所以默认排序用它而不是 created_at
// [EN] This entity has sort_order, the real business default, so it orders by that
var <% .SecondService %>SortableFields = []string{"sort_order", "created_at", "updated_at", "name", "code"}

// <% .SecondService %>Repository <% .SecondService %> 仓库
type <% .SecondService %>Repository struct {
	*CachedRepository[models.<% .SecondService %>Model]
}

// New<% .SecondService %>Repository 创建仓库实例
func New<% .SecondService %>Repository() *<% .SecondService %>Repository {
	return &<% .SecondService %>Repository{
		CachedRepository: NewCachedRepository[models.<% .SecondService %>Model](
			"<% .SecondService | snake %>",
			models.<% .SecondService %>Model{}.TableName(),
		),
	}
}

// ListWithPagination 按请求条件分页查询
func (r *<% .SecondService %>Repository) ListWithPagination(
	ctx context.Context, tenantID string, req *pb.<% .SecondService %>ListRequest,
) ([]*models.<% .SecondService %>Model, *sqlrepo.Pagination32, error) {
	keywordGroup := sqlrepo.NewKeywordFilterGroup(req.GetKeyword().GetValue(), nil, []string{"name", "code"})

	mainGroup := sqlrepo.NewFilterGroup(sqlconstants.LOGIC_AND).
		AddFilter(sqlrepo.NewEqFilter("tenant_id", tenantID)).
		AddEqFilterIfNotEmpty("category", req.GetCategory()).
		AddEqFilterIfNotEmpty("status", req.GetStatus()).
		AddGroupIfNotEmpty(keywordGroup)

	q := sqlrepo.NewQuery().
		WithPaging(int(req.GetPageRequest().GetPage()), int(req.GetPageRequest().GetSize())).
		WithFilterGroup(mainGroup)

	q.AddSafeOrderFromSort(req.GetSort(), "sort_order", "ASC", <% .SecondService %>SortableFields)

	return r.ListWithPagination32(ctx, q)
}
<%- end %>
