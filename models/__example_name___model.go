<%- if eq .GenerateExample "true" %>
/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 00:00:00
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 00:00:00
 * @FilePath: \<% .ProjectName %>\models\<% .ExampleName | snake %>_model.go
 * @Description: <% .ExampleName %> 示例模型定义
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
package models

import (
	"time"
)

// <% .ExampleName %>Model 示例模型 - 字段名称与proto保持一致以支持PBMO自动转换
type <% .ExampleName %>Model struct {
	<% .ExampleName %>Id string    `gorm:"column:<% .ExampleName | snake %>_id;type:varchar(36);primaryKey;comment:<% .ExampleName %>ID" json:"<% .ExampleName | snake %>_id"`
	Name                string    `gorm:"column:name;type:varchar(128);not null;comment:名称" json:"name"`
	Status              int32     `gorm:"column:status;type:tinyint;default:1;not null;comment:状态(1:正常,2:禁用)" json:"status"`
	CreatedAt           time.Time `gorm:"column:created_at;autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at;autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (<% .ExampleName %>Model) TableName() string {
	return "<% .ExampleName | snake %>s"
}

// TableComment 指定表注释
func (<% .ExampleName %>Model) TableComment() string {
	return "<% .ExampleName %>表"
}
<%- end %>
