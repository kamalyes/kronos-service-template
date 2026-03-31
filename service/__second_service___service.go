/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\service\<% .SecondService | snake %>_service.go
 * @Description: <% .SecondService %> 业务逻辑层，实现协议中心的 <% .SecondService %>ServiceServer
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package service
<%- if eq .GenerateExample "true" %>

import (
	"context"

	"github.com/<% .OrgName %>/<% .ProjectName %>/models"
	"github.com/<% .OrgName %>/<% .ProjectName %>/repository"
	pb "<% .ShareProtoModule %>/pb/<% .SecondModule %>"
	commonpb "<% .ShareProtoModule %>/pb/common"
	enumspb "<% .ShareProtoModule %>/pb/enums"
	bizerror "github.com/kamalyes/go-bizerror"
	pbmo "github.com/kamalyes/go-pbmo"
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
)

// <% .SecondService %>Service 实现协议中心 pb.<% .SecondService %>ServiceServer
//
// 与 <% .ExampleService %>Service 的差别值得对照着看：
//   - 更新请求的可选字段全是 wrapper 类型（StringValue / Int32Value），
//     所以部分更新整条交给 pbmo 的 SetIfNotNil / SetInt32Val，一个 if 都不用写
//   - UpdateStatus 与 Delete 返回 common.GeneralEmptyResponse 而不是专属空响应，
//     协议中心对「无返回内容」的两种写法都有示例
//
// [EN] Worth comparing against the other service: here every optional update field is a
// wrapper type, so partial updates defer entirely to pbmo's SetIfNotNil / SetInt32Val with no
// hand-written conditionals; and two RPCs return common.GeneralEmptyResponse rather than
// dedicated empty messages, which the hub demonstrates both ways.
type <% .SecondService %>Service struct {
	pb.Unimplemented<% .SecondService %>ServiceServer
	repo *repository.<% .SecondService %>Repository
}

// New<% .SecondService %>Service 创建服务实例
func New<% .SecondService %>Service(repo *repository.<% .SecondService %>Repository) *<% .SecondService %>Service {
	return &<% .SecondService %>Service{repo: repo}
}

// <% .SecondService %>Create 创建
func (s *<% .SecondService %>Service) <% .SecondService %>Create(
	ctx context.Context, req *pb.<% .SecondService %>CreateRequest,
) (*pb.<% .SecondService %>CreateResponse, error) {
	tenantID, err := tenantOf(ctx)
	if err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsByField(ctx, tenantID, "code", req.GetCode(), "")
	if err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "查询编码是否重复失败", logKV(ctx, "code", req.GetCode(), "error", err)...)
		return nil, dbError(ctx)
	}
	if exists {
		return nil, bizerror.ToGRPCErrorWithTemplate(ctx, bizerror.BizErrCodeConflict,
			map[string]any{"reason": req.GetCode()})
	}

	model := &models.<% .SecondService %>Model{
		Id:        models.NewSnowflakeID(),
		TenantId:  tenantID,
		Name:      req.GetName(),
		Code:      req.GetCode(),
		IconUrl:   req.GetIconUrl(),
		Category:  int32(req.GetCategory()),
		Status:    int32(enumspb.CommonStatus_COMMON_STATUS_NORMAL),
		SortOrder: req.GetSortOrder(),
	}

	if _, err := s.repo.Create(ctx, model); err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "创建失败", logKV(ctx, "code", req.GetCode(), "error", err)...)
		return nil, dbError(ctx)
	}

	info, err := model.ToPB(ctx)
	if err != nil {
		return nil, err
	}
	gwglobal.LOGGER.InfoContextKV(ctx, "创建成功", logKV(ctx, "id", model.Id, "code", model.Code)...)
	return &pb.<% .SecondService %>CreateResponse{Info: info}, nil
}

// <% .SecondService %>Get 查询详情
func (s *<% .SecondService %>Service) <% .SecondService %>Get(
	ctx context.Context, req *pb.<% .SecondService %>GetRequest,
) (*pb.<% .SecondService %>GetResponse, error) {
	tenantID, err := tenantOf(ctx)
	if err != nil {
		return nil, err
	}

	model, err := s.repo.GetByIDCached(ctx, tenantID, req.Get<% .SecondService %>Id())
	if err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "查询失败", logKV(ctx, "id", req.Get<% .SecondService %>Id(), "error", err)...)
		return nil, dbError(ctx)
	}
	if model == nil {
		return nil, notFound(ctx)
	}

	info, err := model.ToPB(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.<% .SecondService %>GetResponse{Info: info}, nil
}

// <% .SecondService %>List 分页列表
func (s *<% .SecondService %>Service) <% .SecondService %>List(
	ctx context.Context, req *pb.<% .SecondService %>ListRequest,
) (*pb.<% .SecondService %>ListResponse, error) {
	tenantID, err := tenantOf(ctx)
	if err != nil {
		return nil, err
	}

	items, paging, err := s.repo.ListWithPagination(ctx, tenantID, req)
	if err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "分页查询失败", logKV(ctx, "error", err)...)
		return nil, dbError(ctx)
	}

	infos, err := models.ToPBs<% .SecondService %>(ctx, items)
	if err != nil {
		return nil, err
	}

	return &pb.<% .SecondService %>ListResponse{
		PageResponse: models.BuildPagination(paging.Page, paging.PageSize, int32(paging.Total)),
		Infos:        infos,
	}, nil
}

// <% .SecondService %>Update 更新
//
// 整条更新由 pbmo 的 SetIfNotNil / SetInt32Val 组装：请求里的可选字段都是 wrapper，
// 「没传」与「传了零值」在类型层面就分得开，不需要任何额外判断。
// SetInt32Val 专收 *wrapperspb.Int32Value，正是 sort_order 的类型。
//
// [EN] Assembled entirely by pbmo: every optional field is a wrapper, so "not sent" and "sent
// zero" are distinguishable at the type level with no extra checks. SetInt32Val takes a
// *wrapperspb.Int32Value, exactly sort_order's type.
func (s *<% .SecondService %>Service) <% .SecondService %>Update(
	ctx context.Context, req *pb.<% .SecondService %>UpdateRequest,
) (*pb.<% .SecondService %>UpdateResponse, error) {
	tenantID, err := tenantOf(ctx)
	if err != nil {
		return nil, err
	}

	id := req.Get<% .SecondService %>Id()
	updates := pbmo.NewUpdates().
		SetIfNotNil("name", req.GetName()).
		SetIfNotNil("icon_url", req.GetIconUrl()).
		SetInt32Val("sort_order", req.GetSortOrder()).
		SetIfNotZero("category", int32(req.GetCategory())).
		Build()

	if err := s.repo.UpdateScoped(ctx, tenantID, id, updates); err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "更新失败", logKV(ctx, "id", id, "error", err)...)
		return nil, dbError(ctx)
	}
	if len(updates) > 0 {
		gwglobal.LOGGER.InfoContextKV(ctx, "更新成功", logKV(ctx, "id", id, "fields", len(updates))...)
	}

	fresh, err := s.repo.GetByIDScoped(ctx, tenantID, id)
	if err != nil {
		return nil, dbError(ctx)
	}
	if fresh == nil {
		return nil, notFound(ctx)
	}
	info, err := fresh.ToPB(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.<% .SecondService %>UpdateResponse{Info: info}, nil
}

// <% .SecondService %>UpdateStatus 更新状态
func (s *<% .SecondService %>Service) <% .SecondService %>UpdateStatus(
	ctx context.Context, req *pb.<% .SecondService %>UpdateStatusRequest,
) (*commonpb.GeneralEmptyResponse, error) {
	tenantID, err := tenantOf(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetStatus() == enumspb.CommonStatus_COMMON_STATUS_UNSPECIFIED {
		return nil, bizerror.ToGRPCErrorWithTemplate(ctx, bizerror.BizErrCodeInvalidParameter,
			map[string]any{"reason": "status"})
	}

	id := req.Get<% .SecondService %>Id()
	updates := pbmo.NewUpdates().Set("status", int32(req.GetStatus())).Build()
	if err := s.repo.UpdateScoped(ctx, tenantID, id, updates); err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "更新状态失败", logKV(ctx, "id", id, "error", err)...)
		return nil, dbError(ctx)
	}
	gwglobal.LOGGER.InfoContextKV(ctx, "更新状态成功", logKV(ctx, "id", id, "status", req.GetStatus().String())...)
	return &commonpb.GeneralEmptyResponse{}, nil
}

// <% .SecondService %>Delete 删除
func (s *<% .SecondService %>Service) <% .SecondService %>Delete(
	ctx context.Context, req *pb.<% .SecondService %>DeleteRequest,
) (*commonpb.GeneralEmptyResponse, error) {
	tenantID, err := tenantOf(ctx)
	if err != nil {
		return nil, err
	}

	id := req.Get<% .SecondService %>Id()
	if err := s.repo.DeleteScoped(ctx, tenantID, id); err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "删除失败", logKV(ctx, "id", id, "error", err)...)
		return nil, dbError(ctx)
	}
	gwglobal.LOGGER.InfoContextKV(ctx, "删除成功", logKV(ctx, "id", id)...)
	return &commonpb.GeneralEmptyResponse{}, nil
}
<%- end %>
