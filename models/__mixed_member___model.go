/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\models\<% .MixedService | snake %>_member_model.go
 * @Description: <% .MixedService %>Member 成员模型
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package models
<%- if eq .GenerateExample "true" %>

import (
	"context"
	"time"

	commonpb "<% .ShareProtoModule %>/pb/common"
	enumspb "<% .ShareProtoModule %>/pb/enums"
	bizerror "github.com/kamalyes/go-bizerror"
	sqlbuilder "github.com/kamalyes/go-sqlbuilder"
)

// <% .MixedService %>MemberRoles 角色集合，落库为 JSON 数组
//
// 角色是「一个成员多个角色」的重复枚举，有三种落库方式：
//   另开一张 member_roles 关联表 —— 需要按角色反查成员时才值得
//   每个角色一列的布尔位 —— 加角色就要改表结构，不可取
//   JSON 数组 —— 本例选它：角色集合总是跟着成员整体读写，没有单独寻址的需求
//
// [EN] Roles are a repeated enum. A join table is worth it only when members must be looked up
// by role; a boolean column per role means a schema change per new role. A JSON array is used
// here because the set is always read and written together with the member.
type <% .MixedService %>MemberRoles = sqlbuilder.Slice[int32]

// <% .MixedService %>MemberModel 对应协议中心 common.<% .MixedService %>MemberInfo
type <% .MixedService %>MemberModel struct {
	Id          string                        `gorm:"column:id;type:varchar(36);primaryKey;comment:成员ID" json:"id"`
	TenantId    string                        `gorm:"column:tenant_id;type:varchar(36);not null;index:idx_<% .MixedService | snake %>_members_tenant;comment:所属<% .MixedService %>ID" json:"tenant_id"`
	UserId      string                        `gorm:"column:user_id;type:varchar(36);not null;comment:用户ID" json:"user_id"`
	DisplayName string                        `gorm:"column:display_name;type:varchar(64);comment:显示名" json:"display_name"`
	Email       string                        `gorm:"column:email;type:varchar(128);comment:邮箱" json:"email"`
	Roles       <% .MixedService %>MemberRoles `gorm:"column:roles;type:json;comment:角色列表" json:"roles"`
	Status      int32                         `gorm:"column:status;type:int;default:0;not null;comment:状态" json:"status"`
	JoinedAt    time.Time                     `gorm:"column:joined_at;autoCreateTime;comment:加入时间" json:"joined_at"`
}

// TableName 指定表名
func (<% .MixedService %>MemberModel) TableName() string {
	return "<% .MixedService | snake %>_members"
}

// TableComment 指定表注释
func (<% .MixedService %>MemberModel) TableComment() string {
	return "<% .MixedService %>成员表"
}

// GetID 实现 repository.TenantScoped
func (m *<% .MixedService %>MemberModel) GetID() string { return m.Id }

// GetTenantID 实现 repository.TenantScoped
func (m *<% .MixedService %>MemberModel) GetTenantID() string { return m.TenantId }

// SetRoles 写入角色集合，入参是 PB 的枚举切片
func (m *<% .MixedService %>MemberModel) SetRoles(roles []enumspb.<% .MixedService %>MemberRole) {
	out := make(<% .MixedService %>MemberRoles, 0, len(roles))
	for _, r := range roles {
		out = append(out, int32(r))
	}
	m.Roles = out
}

// ToPB 转换为协议中心的成员消息
//
// Roles 两侧不同构（落库是 []int32，PB 是 []enums.<% .MixedService %>MemberRole），
// pbmo 配不上，显式转一遍
// [EN] Roles is not structurally identical across the boundary, so it is converted explicitly
func (m *<% .MixedService %>MemberModel) ToPB(ctx context.Context) (*commonpb.<% .MixedService %>MemberInfo, error) {
	info, err := bizerror.ToPB[<% .MixedService %>MemberModel, commonpb.<% .MixedService %>MemberInfo](ctx, m)
	if err != nil {
		return nil, err
	}
	roles := make([]enumspb.<% .MixedService %>MemberRole, 0, len(m.Roles))
	for _, r := range m.Roles {
		roles = append(roles, enumspb.<% .MixedService %>MemberRole(r))
	}
	info.Roles = roles
	return info, nil
}

// ToPBs<% .MixedService %>Member 批量转换
func ToPBs<% .MixedService %>Member(
	ctx context.Context, items []*<% .MixedService %>MemberModel,
) ([]*commonpb.<% .MixedService %>MemberInfo, error) {
	out := make([]*commonpb.<% .MixedService %>MemberInfo, 0, len(items))
	for _, it := range items {
		info, err := it.ToPB(ctx)
		if err != nil {
			return nil, err
		}
		out = append(out, info)
	}
	return out, nil
}
<%- end %>
