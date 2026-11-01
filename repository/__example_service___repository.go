/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 17:05:08
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 19:28:30
 * @FilePath: \<% .ProjectName %>\repository\<% .ExampleService | snake %>_repository.go
 * @Description: <% .ExampleService %> 数据访问层
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package repository
<%- if eq .GenerateExample "true" %>

import (
	"context"

	"github.com/<% .OrgName %>/<% .ProjectName %>/models"
	pb "<% .ShareProtoModule %>/pb/<% .ExampleModule %>"
	sqlconstants "github.com/kamalyes/go-sqlbuilder/constants"
	sqlrepo "github.com/kamalyes/go-sqlbuilder/repository"
)

// <% .ExampleService %>SortableFields 允许排序的字段白名单
//
// 排序字段来自请求（common.Sorting.field），直接拼进 SQL 等于把 ORDER BY 交给调用方。
// 白名单之外的字段由 AddSafeOrderFromSort 丢弃并退回默认排序，而不是报错——
// 调用方传了个不存在的字段时，返回一页数据比返回 500 更合理
//
// [EN] The sort field comes from the request, so splicing it into SQL would hand ORDER BY to
// the caller. AddSafeOrderFromSort drops anything outside this list and falls back to the
// default order rather than erroring: a bogus field is better answered with a page of data
// than a 500
var <% .ExampleService %>SortableFields = []string{"created_at", "updated_at", "name", "code"}

// <% .ExampleService %>Repository <% .ExampleService %> 仓库
//
// 缓存、事务批量删除、租户归属校验都在 CachedRepository 里，这里只放
// 本实体特有的查询（分页条件、唯一性字段）
// [EN] Caching, transactional batch delete and tenant scoping live in CachedRepository; only
// entity-specific queries belong here
type <% .ExampleService %>Repository struct {
	*CachedRepository[models.<% .ExampleService %>Model]
}

// New<% .ExampleService %>Repository 创建仓库实例
func New<% .ExampleService %>Repository() *<% .ExampleService %>Repository {
	return &<% .ExampleService %>Repository{
		CachedRepository: NewCachedRepository[models.<% .ExampleService %>Model](
			"<% .ExampleService | snake %>",
			models.<% .ExampleService %>Model{}.TableName(),
		),
	}
}

// ListWithPagination 按请求条件分页查询
//
// 过滤条件直接由协议中心的 <% .ExampleService %>ListRequest 驱动：
//   - keyword 走 name / code 的 LIKE 组（NewKeywordFilterGroup）
//   - type / status 为零值时不参与过滤（AddEqFilterIfNotEmpty），
//     这与 proto3 的语义一致：枚举零值是 UNSPECIFIED，表示「不限」
//   - 排序过白名单，默认 created_at DESC
//
// [EN] Filters are driven straight by the hub's list request. A zero enum means UNSPECIFIED,
// i.e. "any", so those fields are only applied when non-zero, matching proto3 semantics
func (r *<% .ExampleService %>Repository) ListWithPagination(
	ctx context.Context,
	tenantID string,
	req *pb.<% .ExampleService %>ListRequest,
) ([]*models.<% .ExampleService %>Model, *sqlrepo.Pagination32, error) {
	keywordGroup := sqlrepo.NewKeywordFilterGroup(req.GetKeyword().GetValue(), nil, []string{"name", "code"})

	mainGroup := sqlrepo.NewFilterGroup(sqlconstants.LOGIC_AND).
		AddFilter(sqlrepo.NewEqFilter("tenant_id", tenantID)).
		AddEqFilterIfNotEmpty("type", req.GetType()).
		AddEqFilterIfNotEmpty("status", req.GetStatus()).
		AddGroupIfNotEmpty(keywordGroup)

	q := sqlrepo.NewQuery().
		WithPaging(int(req.GetPageRequest().GetPage()), int(req.GetPageRequest().GetSize())).
		WithFilterGroup(mainGroup)

	q.AddSafeOrderFromSort(req.GetSort(), "created_at", "DESC", <% .ExampleService %>SortableFields)

	return r.ListWithPagination32(ctx, q)
}
<%- end %>
