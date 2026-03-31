/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\models\<% .SecondService | snake %>_model.go
 * @Description: <% .SecondService %> 数据模型，字段与协议中心的 <% .SecondService %>Info 对齐
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package models
<%- if eq .GenerateExample "true" %>

import (
	"context"
	"time"

	commonpb "<% .ShareProtoModule %>/pb/common"
	bizerror "github.com/kamalyes/go-bizerror"
)

// <% .SecondService %>Model 对应协议中心 common.<% .SecondService %>Info
//
// 与 <% .ExampleService %>Model 的区别：这个实体没有嵌套消息与 JSON 列，
// 两侧字段全部同构，所以 ToPB 整条交给 pbmo，不需要任何显式补字段——
// 这是大多数实体的常态，<% .ExampleService %> 那种才是需要额外处理的少数。
//
// [EN] Unlike <% .ExampleService %>Model this entity has no nested messages or JSON columns, so
// every field is structurally identical across the boundary and ToPB defers entirely to pbmo
// with nothing to fill in by hand. That is the common case; the other model is the exception.
type <% .SecondService %>Model struct {
	Id        string    `gorm:"column:id;type:varchar(36);primaryKey;comment:主键ID" json:"id"`
	TenantId  string    `gorm:"column:tenant_id;type:varchar(36);not null;index:idx_<% .SecondService | snake %>s_tenant;comment:租户ID" json:"tenant_id"`
	Name      string    `gorm:"column:name;type:varchar(64);not null;comment:名称" json:"name"`
	Code      string    `gorm:"column:code;type:varchar(32);not null;comment:编码" json:"code"`
	IconUrl   string    `gorm:"column:icon_url;type:varchar(512);comment:图标地址" json:"icon_url"`
	Category  int32     `gorm:"column:category;type:int;default:0;not null;comment:分类" json:"category"`
	Status    int32     `gorm:"column:status;type:int;default:0;not null;index:idx_<% .SecondService | snake %>s_status;comment:状态" json:"status"`
	SortOrder int32     `gorm:"column:sort_order;type:int;default:0;not null;comment:排序序号" json:"sort_order"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (<% .SecondService %>Model) TableName() string {
	return "<% .SecondService | snake %>s"
}

// TableComment 指定表注释
func (<% .SecondService %>Model) TableComment() string {
	return "<% .SecondService %>表"
}

// GetID 实现 repository.TenantScoped
func (m *<% .SecondService %>Model) GetID() string { return m.Id }

// GetTenantID 实现 repository.TenantScoped
func (m *<% .SecondService %>Model) GetTenantID() string { return m.TenantId }

// ToPB 转换为协议中心的 <% .SecondService %>Info
func (m *<% .SecondService %>Model) ToPB(ctx context.Context) (*commonpb.<% .SecondService %>Info, error) {
	return bizerror.ToPB[<% .SecondService %>Model, commonpb.<% .SecondService %>Info](ctx, m)
}

// ToPBs<% .SecondService %> 批量转换
//
// 直接用 bizerror.ToPBs：整个切片一次转换，内部复用同一份字段计划，
// 不像逐条调用那样每次都走一遍类型断言
// [EN] bizerror.ToPBs converts the whole slice reusing one field plan
func ToPBs<% .SecondService %>(ctx context.Context, items []*<% .SecondService %>Model) ([]*commonpb.<% .SecondService %>Info, error) {
	return bizerror.ToPBs[<% .SecondService %>Model, commonpb.<% .SecondService %>Info](ctx, items)
}
<%- end %>
