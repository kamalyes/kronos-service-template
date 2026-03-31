/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\service\<% .MixedService | snake %>_member_service.go
 * @Description: <% .MixedService %>Member 业务逻辑，实现协议中心的 <% .MixedService %>MemberServiceServer
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package service
<%- if eq .GenerateExample "true" %>

import (
	"context"

	"github.com/<% .OrgName %>/<% .ProjectName %>/models"
	"github.com/<% .OrgName %>/<% .ProjectName %>/repository"
	pb "<% .ShareProtoModule %>/pb/<% .MixedModule %>"
	commonpb "<% .ShareProtoModule %>/pb/common"
	enumspb "<% .ShareProtoModule %>/pb/enums"
	bizerror "github.com/kamalyes/go-bizerror"
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
)

// <% .MixedService %>MemberService 实现协议中心 pb.<% .MixedService %>MemberServiceServer
//
// 与 <% .MixedService %>Service 同处 <% .MixedModule %> 模块、共用 proto 包与 Go 包，
// 但两者不互相 import：共享的类型在 common/ 里。协议中心的规范第 8.5 节写明
// 「拆开后两个 service 不要互相 import，否则又被焊回一起」。
//
// 这里的父 ID 来自请求字段（<% .MixedService | snake %>_id）而不是 gwmiddleware：
// 成员接口操作的是「某个聚合根下的成员」，父 ID 是路径参数，属于寻址信息；
// 真实系统里还应校验调用者对该聚合根有权限，那是鉴权层的事，不在这里。
//
// [EN] It shares the module and Go package with the aggregate service but never imports it;
// shared types live in common/. The parent id comes from the request field rather than the
// context because it is addressing information, a path parameter. A real system would also
// verify the caller's permission on that root, which belongs to the authorization layer.
type <% .MixedService %>MemberService struct {
	pb.Unimplemented<% .MixedService %>MemberServiceServer
	repo *repository.<% .MixedService %>MemberRepository
}

// New<% .MixedService %>MemberService 创建服务实例
func New<% .MixedService %>MemberService(repo *repository.<% .MixedService %>MemberRepository) *<% .MixedService %>MemberService {
	return &<% .MixedService %>MemberService{repo: repo}
}

// <% .MixedService %>MemberList 分页列出成员
//
// 子资源列表同样要分页：它会随父资源增长，不分页的「列出全部成员」
// 在大租户上就是一次全表扫描
// [EN] A sub-resource list is paged too: it grows with the parent, and listing all members
// unpaged becomes a full scan on a large tenant
func (s *<% .MixedService %>MemberService) <% .MixedService %>MemberList(
	ctx context.Context, req *pb.<% .MixedService %>MemberListRequest,
) (*pb.<% .MixedService %>MemberListResponse, error) {
	tenantID := req.Get<% .MixedService %>Id()

	items, paging, err := s.repo.ListWithPagination(ctx, tenantID, req)
	if err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "分页查询成员失败", logKV(ctx, "id", tenantID, "error", err)...)
		return nil, dbError(ctx)
	}

	infos, err := models.ToPBs<% .MixedService %>Member(ctx, items)
	if err != nil {
		return nil, err
	}

	return &pb.<% .MixedService %>MemberListResponse{
		PageResponse: models.BuildPagination(paging.Page, paging.PageSize, int32(paging.Total)),
		Infos:        infos,
	}, nil
}

// <% .MixedService %>MemberAdd 添加成员
func (s *<% .MixedService %>MemberService) <% .MixedService %>MemberAdd(
	ctx context.Context, req *pb.<% .MixedService %>MemberAddRequest,
) (*pb.<% .MixedService %>MemberAddResponse, error) {
	tenantID := req.Get<% .MixedService %>Id()

	exists, err := s.repo.UserExists(ctx, tenantID, req.GetUserId())
	if err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "查询成员是否已存在失败",
			logKV(ctx, "id", tenantID, "user_id", req.GetUserId(), "error", err)...)
		return nil, dbError(ctx)
	}
	if exists {
		return nil, bizerror.ToGRPCErrorWithTemplate(ctx, bizerror.BizErrCodeConflict,
			map[string]any{"reason": req.GetUserId()})
	}

	member := &models.<% .MixedService %>MemberModel{
		Id:       models.NewSnowflakeID(),
		TenantId: tenantID,
		UserId:   req.GetUserId(),
		Email:    req.GetEmail(),
		Status:   int32(enumspb.CommonStatus_COMMON_STATUS_NORMAL),
	}
	member.SetRoles(req.GetRoles())

	if _, err := s.repo.Create(ctx, member); err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "添加成员失败", logKV(ctx, "id", tenantID, "error", err)...)
		return nil, dbError(ctx)
	}

	info, err := member.ToPB(ctx)
	if err != nil {
		return nil, err
	}
	gwglobal.LOGGER.InfoContextKV(ctx, "添加成员成功",
		logKV(ctx, "id", tenantID, "member_id", member.Id, "roles", len(req.GetRoles()))...)
	return &pb.<% .MixedService %>MemberAddResponse{Info: info}, nil
}

// <% .MixedService %>MemberRemove 移除成员
func (s *<% .MixedService %>MemberService) <% .MixedService %>MemberRemove(
	ctx context.Context, req *pb.<% .MixedService %>MemberRemoveRequest,
) (*commonpb.GeneralEmptyResponse, error) {
	tenantID := req.Get<% .MixedService %>Id()
	memberID := req.GetMemberId()

	if err := s.repo.DeleteScoped(ctx, tenantID, memberID); err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "移除成员失败",
			logKV(ctx, "id", tenantID, "member_id", memberID, "error", err)...)
		return nil, dbError(ctx)
	}
	gwglobal.LOGGER.InfoContextKV(ctx, "移除成员成功",
		logKV(ctx, "id", tenantID, "member_id", memberID)...)
	return &commonpb.GeneralEmptyResponse{}, nil
}

// <% .MixedService %>MemberRoleReplace 覆盖成员角色
//
// 语义是覆盖而不是追加，协议中心的字段注释写的是「覆盖后的角色集合」。
// 覆盖型接口不能用「空集合视为不变」那套：调用方传空集合就是要清空角色，
// 而「不想改角色」就根本不该调这个接口
//
// [EN] The semantics are replace, not append, per the hub's field comment. A replace endpoint
// must not treat an empty set as "unchanged": an empty set means clear the roles, and a caller
// who does not want to change them simply does not call this endpoint
func (s *<% .MixedService %>MemberService) <% .MixedService %>MemberRoleReplace(
	ctx context.Context, req *pb.<% .MixedService %>MemberRoleReplaceRequest,
) (*pb.<% .MixedService %>MemberRoleReplaceResponse, error) {
	tenantID := req.Get<% .MixedService %>Id()
	memberID := req.GetMemberId()

	member, err := s.repo.GetByIDScoped(ctx, tenantID, memberID)
	if err != nil {
		return nil, dbError(ctx)
	}
	if member == nil {
		return nil, notFound(ctx)
	}

	member.SetRoles(req.GetRoles())
	if err := s.repo.UpdateScoped(ctx, tenantID, memberID,
		map[string]interface{}{"roles": member.Roles}); err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "覆盖成员角色失败",
			logKV(ctx, "member_id", memberID, "error", err)...)
		return nil, dbError(ctx)
	}

	info, err := member.ToPB(ctx)
	if err != nil {
		return nil, err
	}
	gwglobal.LOGGER.InfoContextKV(ctx, "覆盖成员角色成功",
		logKV(ctx, "member_id", memberID, "roles", len(req.GetRoles()))...)
	return &pb.<% .MixedService %>MemberRoleReplaceResponse{Info: info}, nil
}

// <% .MixedService %>MemberBatchInvite 批量邀请
//
// 批量接口的三条约定（协议中心规范第 8.4 节）：
//   走 POST 而不是带 body 的 DELETE/PUT——网关与代理对带 body 的 DELETE 行为不一致
//   必须有条数上限，由 proto 的 validate:"required,max=50" 保证，这里不重复校验
//   逐条返回结果（common.BatchResult），让调用方精确重试，而不是整体成功/失败
//
// 已存在的邮箱一次查完再比对，不是逐个查：50 条逐个查就是 50 次往返。
// 真正要插入的那批放进一个事务——部分成功后调用方重试会把已建的再建一遍，
// 而成员表上「同租户同邮箱」没有唯一约束（邮箱可为空），重试就产生重复行。
//
// [EN] The three conventions for batch endpoints: POST rather than a body-bearing DELETE, since
// gateways and proxies treat those inconsistently; a mandatory count cap, enforced by the
// proto's validate tag and not re-checked here; and per-item results so the caller retries
// precisely. Existing emails are fetched in one query rather than one per item, and the rows to
// insert share a transaction because a partial success plus a retry would duplicate rows: there
// is no unique constraint on tenant plus email, as email may be empty.
func (s *<% .MixedService %>MemberService) <% .MixedService %>MemberBatchInvite(
	ctx context.Context, req *pb.<% .MixedService %>MemberBatchInviteRequest,
) (*pb.<% .MixedService %>MemberBatchInviteResponse, error) {
	tenantID := req.Get<% .MixedService %>Id()
	emails := req.GetEmails()

	existing, err := s.repo.EmailsExist(ctx, tenantID, emails)
	if err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "查询已存在邮箱失败", logKV(ctx, "id", tenantID, "error", err)...)
		return nil, dbError(ctx)
	}

	results := make([]*commonpb.BatchResult, 0, len(emails))
	toCreate := make([]*models.<% .MixedService %>MemberModel, 0, len(emails))
	seen := make(map[string]bool, len(emails))

	for _, email := range emails {
		// 入参自身可能重复，去重后再判断，否则同一邮箱会被插两次
		// [EN] The input itself may repeat; dedupe first or the same email is inserted twice
		if seen[email] {
			results = append(results, &commonpb.BatchResult{
				Id:           email,
				Success:      false,
				ErrorMessage: gwmiddlewareMsg(ctx, bizerror.BizErrCodeConflict),
			})
			continue
		}
		seen[email] = true

		if existing[email] {
			results = append(results, &commonpb.BatchResult{
				Id:           email,
				Success:      false,
				ErrorMessage: gwmiddlewareMsg(ctx, bizerror.BizErrCodeConflict),
			})
			continue
		}

		member := &models.<% .MixedService %>MemberModel{
			Id:       models.NewSnowflakeID(),
			TenantId: tenantID,
			Email:    email,
			// CommonStatus 只有 UNSPECIFIED / NORMAL / DISABLED，没有 PENDING：
			// 被邀请但未接受的成员先置 DISABLED，接受后再转 NORMAL
			// [EN] CommonStatus has no PENDING, so an invited-but-unaccepted member starts
			// DISABLED and moves to NORMAL on acceptance
			Status:   int32(enumspb.CommonStatus_COMMON_STATUS_DISABLED),
		}
		member.SetRoles(req.GetRoles())
		toCreate = append(toCreate, member)
		results = append(results, &commonpb.BatchResult{Id: email, Success: true})
	}

	if err := s.repo.CreateBatchInTx(ctx, toCreate); err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "批量邀请失败",
			logKV(ctx, "id", tenantID, "count", len(toCreate), "error", err)...)
		return nil, dbError(ctx)
	}

	gwglobal.LOGGER.InfoContextKV(ctx, "批量邀请完成",
		logKV(ctx, "id", tenantID, "total", len(emails), "created", len(toCreate))...)
	return &pb.<% .MixedService %>MemberBatchInviteResponse{Results: results}, nil
}
<%- end %>
