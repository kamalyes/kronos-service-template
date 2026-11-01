/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 17:05:08
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 19:28:30
 * @FilePath: \<% .ProjectName %>\repository\<% .MixedService | snake %>_member_repository.go
 * @Description: <% .MixedService %>Member 数据访问层
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package repository
<%- if eq .GenerateExample "true" %>

import (
	"context"

	"github.com/<% .OrgName %>/<% .ProjectName %>/models"
	pb "<% .ShareProtoModule %>/pb/<% .MixedModule %>"
	sqlconstants "github.com/kamalyes/go-sqlbuilder/constants"
	sqlrepo "github.com/kamalyes/go-sqlbuilder/repository"
)

// <% .MixedService %>MemberRepository 成员仓库
//
// 成员是租户内实体（有 tenant_id），所以走泛型基座，缓存与事务批量删除都是现成的
// [EN] Members are tenant-scoped, so the generic base applies and caching plus transactional
// batch delete come for free
type <% .MixedService %>MemberRepository struct {
	*CachedRepository[models.<% .MixedService %>MemberModel]
}

// New<% .MixedService %>MemberRepository 创建仓库实例
func New<% .MixedService %>MemberRepository() *<% .MixedService %>MemberRepository {
	return &<% .MixedService %>MemberRepository{
		CachedRepository: NewCachedRepository[models.<% .MixedService %>MemberModel](
			"<% .MixedService | snake %>_member",
			models.<% .MixedService %>MemberModel{}.TableName(),
		),
	}
}

// ListWithPagination 分页查询成员
//
// 按角色过滤要查 JSON 列：角色落库是 JSON 数组，用 OP_JSON_CONTAINS 而不是等值比较。
// 这也是选 JSON 落库的代价——换成关联表的话这里是一次 JOIN，索引更好用；
// 本例按角色反查不是主路径，所以接受这个代价
//
// [EN] Filtering by role queries the JSON column with a contains operator rather than equality.
// This is the cost of the JSON choice: a join table would index better here, but looking up by
// role is not the primary path in this example
func (r *<% .MixedService %>MemberRepository) ListWithPagination(
	ctx context.Context, tenantID string, req *pb.<% .MixedService %>MemberListRequest,
) ([]*models.<% .MixedService %>MemberModel, *sqlrepo.Pagination32, error) {
	mainGroup := sqlrepo.NewFilterGroup(sqlconstants.LOGIC_AND).
		AddFilter(sqlrepo.NewEqFilter("tenant_id", tenantID)).
		AddEqFilterIfNotEmpty("status", req.GetStatus())

	if role := req.GetRole(); role != 0 {
		mainGroup.AddFilter(sqlrepo.NewFilter("roles", sqlconstants.OP_JSON_CONTAINS, int32(role)))
	}

	q := sqlrepo.NewQuery().
		WithPaging(int(req.GetPageRequest().GetPage()), int(req.GetPageRequest().GetSize())).
		WithFilterGroup(mainGroup).
		AddOrder("joined_at", "DESC")

	return r.ListWithPagination32(ctx, q)
}

// UserExists 判断用户是否已是该聚合根的成员
func (r *<% .MixedService %>MemberRepository) UserExists(ctx context.Context, tenantID, userID string) (bool, error) {
	count, err := r.Count(ctx,
		sqlrepo.NewEqFilter("tenant_id", tenantID),
		sqlrepo.NewEqFilter("user_id", userID),
	)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// EmailsExist 批量查出已存在的邮箱集合，供批量邀请去重
//
// 一次查完而不是逐个查：批量邀请最多 50 条（协议中心的 validate:"required,max=50"），
// 逐个查就是 50 次往返
// [EN] One query instead of per-item: a batch invite carries up to 50 emails per the hub's
// validate tag, and checking them one by one would be 50 round trips
func (r *<% .MixedService %>MemberRepository) EmailsExist(
	ctx context.Context, tenantID string, emails []string,
) (map[string]bool, error) {
	if len(emails) == 0 {
		return map[string]bool{}, nil
	}
	q := sqlrepo.NewQuery().
		AddFilter(sqlrepo.NewEqFilter("tenant_id", tenantID)).
		AddFilter(sqlrepo.NewFilter("email", sqlconstants.OP_IN, emails))

	items, err := r.List(ctx, q)
	if err != nil {
		return nil, err
	}
	out := make(map[string]bool, len(items))
	for _, it := range items {
		out[it.Email] = true
	}
	return out, nil
}

// CreateBatchInTx 事务内批量创建成员
//
// 批量邀请要么全建要么全不建：部分成功后调用方重试会把已建的那批再建一遍，
// 而成员表上「同一租户同一邮箱」并没有唯一约束（邮箱可为空），重试就会产生重复行
//
// [EN] All or nothing: after a partial success a retry would re-create the ones that succeeded,
// and since there is no unique constraint on tenant plus email (email may be empty) the retry
// would produce duplicate rows
func (r *<% .MixedService %>MemberRepository) CreateBatchInTx(
	ctx context.Context, items []*models.<% .MixedService %>MemberModel,
) error {
	if len(items) == 0 {
		return nil
	}
	return r.Transaction(ctx, func(tx sqlrepo.Transaction[models.<% .MixedService %>MemberModel]) error {
		return tx.CreateBatch(ctx, items...)
	})
}
<%- end %>
