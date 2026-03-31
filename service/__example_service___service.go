/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\service\<% .ExampleService | snake %>_service.go
 * @Description: <% .ExampleService %> 业务逻辑层，实现协议中心的 <% .ExampleService %>ServiceServer
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package service
<%- if eq .GenerateExample "true" %>

import (
	"context"

	"github.com/<% .OrgName %>/<% .ProjectName %>/models"
	"github.com/<% .OrgName %>/<% .ProjectName %>/repository"
	pb "<% .ShareProtoModule %>/pb/<% .ExampleModule %>"
	commonpb "<% .ShareProtoModule %>/pb/common"
	enumspb "<% .ShareProtoModule %>/pb/enums"
	bizerror "github.com/kamalyes/go-bizerror"
	pbmo "github.com/kamalyes/go-pbmo"
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
	gwmiddleware "github.com/kamalyes/go-rpc-gateway/middleware"
	sqltypes "github.com/kamalyes/go-sqlbuilder/types"
)

// <% .ExampleService %>Service 实现协议中心 pb.<% .ExampleService %>ServiceServer
//
// 这是一份可直接改的真实示例，不是占位：它实现的就是协议中心
// proto/<% .ExampleModule %>/<% .ExampleService | snake %>_service.proto 里定义的那个 service，
// 换成自己的业务时照着替换 model 与 RPC 体即可。
//
// 嵌入 Unimplemented 是 protoc-gen-go-grpc 的要求，也给了向前兼容：
// 协议中心新增 RPC 后本服务仍能编译，新方法返回 Unimplemented 而不是编译失败。
//
// [EN] A real, editable example rather than a placeholder: it implements the very service the
// contract hub defines. Embedding Unimplemented is required by protoc-gen-go-grpc and gives
// forward compatibility: a new RPC in the hub still compiles here and returns Unimplemented.
type <% .ExampleService %>Service struct {
	pb.Unimplemented<% .ExampleService %>ServiceServer
	repo *repository.<% .ExampleService %>Repository
}

// New<% .ExampleService %>Service 创建服务实例
func New<% .ExampleService %>Service(repo *repository.<% .ExampleService %>Repository) *<% .ExampleService %>Service {
	return &<% .ExampleService %>Service{repo: repo}
}

// <% .ExampleService %>Create 创建
//
// 参数校验不在这里写：协议中心的 proto 字段上注入了 @inject_tag: validate:"..."，
// go-rpc-gateway 已把基于 go-argus 的 struct tag 校验拦截器挂在链上，
// 走到这里时 req 的格式约束已经满足。这里只做格式之外的业务规则（唯一性之类）。
//
// [EN] Field validation is not written here: the hub injects validate tags into the proto
// fields and go-rpc-gateway runs the go-argus struct-tag interceptor, so format constraints
// already hold. Only business rules beyond format belong here.
func (s *<% .ExampleService %>Service) <% .ExampleService %>Create(
	ctx context.Context, req *pb.<% .ExampleService %>CreateRequest,
) (*pb.<% .ExampleService %>CreateResponse, error) {
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
		// 带 {reason} 模板的错误码：文案由 go-bizerror 按请求语言渲染，
		// 这里只提供原因，不拼人类可读的句子
		// [EN] A templated code: go-bizerror renders the message in the request language
		return nil, bizerror.ToGRPCErrorWithTemplate(ctx, bizerror.BizErrCodeConflict,
			map[string]any{"reason": req.GetCode()})
	}

	model := &models.<% .ExampleService %>Model{
		Id:       models.NewSnowflakeID(),
		TenantId: tenantID,
		Name:     req.GetName(),
		Code:     req.GetCode(),
		Type:     int32(req.GetType()),
		Level:    int32(req.GetLevel()),
		Status:   int32(enumspb.CommonStatus_COMMON_STATUS_NORMAL),
		Tags:     req.GetTags(),
	}
	model.SetContact(req.GetContact())

	if _, err := s.repo.Create(ctx, model); err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "创建失败", logKV(ctx, "code", req.GetCode(), "error", err)...)
		return nil, dbError(ctx)
	}

	info, err := model.ToPB(ctx)
	if err != nil {
		return nil, err
	}
	gwglobal.LOGGER.InfoContextKV(ctx, "创建成功", logKV(ctx, "id", model.Id, "code", model.Code)...)
	return &pb.<% .ExampleService %>CreateResponse{Info: info}, nil
}

// <% .ExampleService %>Get 查询详情
func (s *<% .ExampleService %>Service) <% .ExampleService %>Get(
	ctx context.Context, req *pb.<% .ExampleService %>GetRequest,
) (*pb.<% .ExampleService %>GetResponse, error) {
	tenantID, err := tenantOf(ctx)
	if err != nil {
		return nil, err
	}

	model, err := s.repo.GetByIDCached(ctx, tenantID, req.Get<% .ExampleService %>Id())
	if err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "查询失败", logKV(ctx, "id", req.Get<% .ExampleService %>Id(), "error", err)...)
		return nil, dbError(ctx)
	}
	if model == nil {
		return nil, notFound(ctx)
	}

	info, err := model.ToPB(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.<% .ExampleService %>GetResponse{Info: info}, nil
}

// <% .ExampleService %>List 分页列表
func (s *<% .ExampleService %>Service) <% .ExampleService %>List(
	ctx context.Context, req *pb.<% .ExampleService %>ListRequest,
) (*pb.<% .ExampleService %>ListResponse, error) {
	tenantID, err := tenantOf(ctx)
	if err != nil {
		return nil, err
	}

	items, paging, err := s.repo.ListWithPagination(ctx, tenantID, req)
	if err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "分页查询失败", logKV(ctx, "error", err)...)
		return nil, dbError(ctx)
	}

	infos, err := models.ToPBs<% .ExampleService %>(ctx, items)
	if err != nil {
		return nil, err
	}

	return &pb.<% .ExampleService %>ListResponse{
		// Total 在 PaginationT 里固定是 int64（为支持大数据量），而 common.Paging.total
		// 是 int32，这里显式收窄；真到 21 亿条时分页本身已不是合适的读取方式
		// [EN] Total is int64 in PaginationT while common.Paging.total is int32
		PageResponse: models.BuildPagination(paging.Page, paging.PageSize, int32(paging.Total)),
		Infos:        infos,
	}, nil
}

// <% .ExampleService %>Update 更新
//
// 用 pbmo.NewUpdates() 拼部分更新：请求里的可选字段是 *wrappers.StringValue，
// SetIfNotNil 只在调用方真的传了该字段时才写进 updates。
// 换成「取 GetValue() 再判空串」会把「传了空串」和「没传」混为一谈，
// 于是清空某个字段永远做不到。
//
// [EN] Partial updates are assembled with pbmo.NewUpdates(): optional request fields are
// wrapper types, and SetIfNotNil only writes when the caller actually sent the field. Reading
// GetValue() and testing for "" would conflate "sent empty" with "not sent", making it
// impossible to clear a field.
func (s *<% .ExampleService %>Service) <% .ExampleService %>Update(
	ctx context.Context, req *pb.<% .ExampleService %>UpdateRequest,
) (*pb.<% .ExampleService %>UpdateResponse, error) {
	tenantID, err := tenantOf(ctx)
	if err != nil {
		return nil, err
	}

	id := req.Get<% .ExampleService %>Id()
	model, err := s.repo.GetByIDCached(ctx, tenantID, id)
	if err != nil {
		return nil, dbError(ctx)
	}
	if model == nil {
		return nil, notFound(ctx)
	}

	builder := pbmo.NewUpdates().
		SetIfNotNil("name", req.GetName()).
		SetIfNotZero("level", int32(req.GetLevel()))

	// tags 用一个独立的布尔开关控制是否覆盖：空数组既可能是「清空」也可能是「没传」，
	// 协议中心为此在请求里专门加了 update_tags
	// [EN] An empty array means either "clear" or "not sent", so the hub added update_tags
	if req.GetUpdateTags().GetValue() {
		builder.SetJSONSlice("tags", req.GetTags())
	}
	if req.GetContact() != nil {
		// SetJSON 的入参是「已序列化好的 JSON 字符串」，不是对象；
		// 这里直接 Set 一个 ProtoJSON 包装体，由它的 driver.Valuer 负责序列化，
		// 和 model 落库走的是同一条路径，不会出现两种 JSON 形状
		// [EN] SetJSON takes an already-serialized JSON string, not an object. Setting the
		// ProtoJSON wrapper instead reuses its driver.Valuer, so the stored shape matches
		// exactly what the model writes
		builder.Set("contact", sqltypes.ProtoJSON[*commonpb.<% .ExampleService %>Contact]{Data: req.GetContact()})
	}

	updates := builder.Build()
	if len(updates) > 0 {
		if err := s.repo.UpdateScoped(ctx, tenantID, id, updates); err != nil {
			gwglobal.LOGGER.ErrorContextKV(ctx, "更新失败", logKV(ctx, "id", id, "error", err)...)
			return nil, dbError(ctx)
		}
		gwglobal.LOGGER.InfoContextKV(ctx, "更新成功", logKV(ctx, "id", id, "fields", len(updates))...)
	}

	fresh, err := s.repo.GetByIDCached(ctx, tenantID, id)
	if err != nil || fresh == nil {
		return nil, dbError(ctx)
	}
	info, err := fresh.ToPB(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.<% .ExampleService %>UpdateResponse{Info: info}, nil
}

// <% .ExampleService %>UpdateStatus 更新状态（局部更新）
func (s *<% .ExampleService %>Service) <% .ExampleService %>UpdateStatus(
	ctx context.Context, req *pb.<% .ExampleService %>UpdateStatusRequest,
) (*pb.<% .ExampleService %>UpdateStatusResponse, error) {
	tenantID, err := tenantOf(ctx)
	if err != nil {
		return nil, err
	}

	// 枚举的允许值由业务方法校验，不用 oneof 系列规则：
	// go-argus 会把 protobuf 枚举格式化成名称再比较，数字候选值会误拒合法请求
	// （协议中心的 proto 规范第 7 节写了这条）
	// [EN] Allowed enum values are checked here rather than with oneof-style tags
	if req.GetStatus() == enumspb.CommonStatus_COMMON_STATUS_UNSPECIFIED {
		return nil, bizerror.ToGRPCErrorWithTemplate(ctx, bizerror.BizErrCodeInvalidParameter,
			map[string]any{"reason": "status"})
	}

	id := req.Get<% .ExampleService %>Id()
	// SetInt32Val 收的是 *wrapperspb.Int32Value（给可选字段用）；
	// status 这里是必填枚举，已是确定值，直接 Set
	// [EN] SetInt32Val takes a wrapper for optional fields; status is a required enum here
	updates := pbmo.NewUpdates().Set("status", int32(req.GetStatus())).Build()
	if err := s.repo.UpdateScoped(ctx, tenantID, id, updates); err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "更新状态失败", logKV(ctx, "id", id, "error", err)...)
		return nil, dbError(ctx)
	}
	gwglobal.LOGGER.InfoContextKV(ctx, "更新状态成功", logKV(ctx, "id", id, "status", req.GetStatus().String())...)
	return &pb.<% .ExampleService %>UpdateStatusResponse{}, nil
}

// <% .ExampleService %>Delete 删除
func (s *<% .ExampleService %>Service) <% .ExampleService %>Delete(
	ctx context.Context, req *pb.<% .ExampleService %>DeleteRequest,
) (*pb.<% .ExampleService %>DeleteResponse, error) {
	tenantID, err := tenantOf(ctx)
	if err != nil {
		return nil, err
	}

	id := req.Get<% .ExampleService %>Id()
	if err := s.repo.DeleteScoped(ctx, tenantID, id); err != nil {
		gwglobal.LOGGER.ErrorContextKV(ctx, "删除失败", logKV(ctx, "id", id, "error", err)...)
		return nil, dbError(ctx)
	}
	gwglobal.LOGGER.InfoContextKV(ctx, "删除成功", logKV(ctx, "id", id)...)
	return &pb.<% .ExampleService %>DeleteResponse{}, nil
}

// <% .ExampleService %>BatchDelete 批量删除
//
// 逐条返回结果而不是整体成败：调用方据此只重试失败的那几条。
// 条数上限由 proto 的 validate:"required,max=100" 保证，这里不必再查一遍。
//
// [EN] Per-item results rather than all-or-nothing, so the caller retries only what failed.
// The count cap is enforced by the proto's validate tag, so it is not re-checked here.
func (s *<% .ExampleService %>Service) <% .ExampleService %>BatchDelete(
	ctx context.Context, req *pb.<% .ExampleService %>BatchDeleteRequest,
) (*pb.<% .ExampleService %>BatchDeleteResponse, error) {
	tenantID, err := tenantOf(ctx)
	if err != nil {
		return nil, err
	}

	ids := req.Get<% .ExampleService %>Ids()
	results := make([]*commonpb.BatchResult, 0, len(ids))
	var failed int
	for _, id := range ids {
		if err := s.repo.DeleteScoped(ctx, tenantID, id); err != nil {
			failed++
			results = append(results, &commonpb.BatchResult{
				Id:      id,
				Success: false,
				ErrorMessage: gwmiddleware.GetMsgByKey(ctx, bizerror.BizErrCodeDatabaseOperationFailed),
			})
			continue
		}
		results = append(results, &commonpb.BatchResult{Id: id, Success: true})
	}

	gwglobal.LOGGER.InfoContextKV(ctx, "批量删除完成",
		logKV(ctx, "total", len(ids), "failed", failed)...)
	return &pb.<% .ExampleService %>BatchDeleteResponse{Results: results}, nil
}
<%- end %>
