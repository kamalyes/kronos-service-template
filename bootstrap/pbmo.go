/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\bootstrap\pbmo.go
 * @Description: Model-PB 转换器注册
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package bootstrap

import (
<%- if eq .GenerateExample "true" %>
	"github.com/<% .OrgName %>/<% .ProjectName %>/models"
	commonpb "<% .ShareProtoModule %>/pb/common"
	pbmo "github.com/kamalyes/go-pbmo"
<%- end %>
)

// registerConverters 注册所有 Model-PB 转换器
//
// 必须在服务开始处理请求前注册：pbmo 在 Register 时按反射构建字段搬运计划并缓存，
// 之后每次 ToPB 走的是预编译好的闭包。没注册过的类型对在首次转换时才构建，
// 那一次会落在请求路径上。
//
// 字段名两侧一致即自动配对（不分大小写），名字对不上时用
// .WithFieldMapping("ModelField", "PBField") 显式指定，例如
// pbmo.Register[models.XxxModel, pb.XxxInfo]().WithFieldMapping("ID", "XxxId")
//
// [EN] Register before the service starts serving: pbmo builds and caches a field-copy plan
// via reflection at registration time, and later conversions run the precompiled closure. An
// unregistered pair is built on first use, which lands that cost on a request. Fields pair
// automatically by case-insensitive name; use WithFieldMapping when the names differ.
func registerConverters() {
<%- if eq .GenerateExample "true" %>
	pbmo.Register[models.<% .ExampleService %>Model, commonpb.<% .ExampleService %>Info]()
	pbmo.Register[models.<% .SecondService %>Model, commonpb.<% .SecondService %>Info]()
	pbmo.Register[models.<% .MixedService %>Model, commonpb.<% .MixedService %>Info]()
	pbmo.Register[models.<% .MixedService %>MemberModel, commonpb.<% .MixedService %>MemberInfo]()
	// IsPrimary ↔ Primary 名字对不上，显式声明映射；其余字段同名自动配对
	// [EN] IsPrimary and Primary differ by name, so the mapping is declared explicitly
	pbmo.Register[models.<% .MixedService %>RegionAssignmentModel, commonpb.<% .MixedService %>RegionAssignment]().
		WithFieldMapping("IsPrimary", "Primary")
<%- else %>
	// TODO: 注册本服务的 Model-PB 转换器
	// pbmo.Register[models.XxxModel, commonpb.XxxInfo]()
<%- end %>
}
