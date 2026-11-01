/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 17:05:08
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 19:28:30
 * @FilePath: \<% .ProjectName %>\models\table_names.go
 * @Description: 数据库表名常量定义
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */

package models

// 数据库表名常量 - 通过调用模型的 TableName() 方法动态获取
var (
<%- if eq .GenerateExample "true" %>
	TableName<% .ExampleService %> = <% .ExampleService %>Model{}.TableName()
	TableName<% .SecondService %>  = <% .SecondService %>Model{}.TableName()
	TableName<% .MixedService %>   = <% .MixedService %>Model{}.TableName()
	TableName<% .MixedService %>RegionAssignment = <% .MixedService %>RegionAssignmentModel{}.TableName()
	TableName<% .MixedService %>Member = <% .MixedService %>MemberModel{}.TableName()
<%- else %>
	// TODO: 添加表名常量
	// TableNameExample = ExampleModel{}.TableName()
<%- end %>
)
