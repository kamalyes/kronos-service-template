/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 17:05:08
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 19:28:30
 * @FilePath: \<% .ProjectName %>\migrations\auto_migrate.go
 * @Description: 数据库自动迁移 基于 go-sqlbuilder Migrator
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */

package migrations

import (
<%- if eq .GenerateExample "true" %>
	"github.com/<% .OrgName %>/<% .ProjectName %>/models"
<%- end %>
	"github.com/kamalyes/go-sqlbuilder/db"
	"gorm.io/gorm"
)

// GetAllModels 获取所有需要初始化的数据库模型
func GetAllModels() []interface{} {
	return []interface{}{
<%- if eq .GenerateExample "true" %>
		&models.<% .ExampleService %>Model{},
		&models.<% .SecondService %>Model{},
		&models.<% .MixedService %>Model{},
		&models.<% .MixedService %>RegionAssignmentModel{},
		&models.<% .MixedService %>MemberModel{},
<%- else %>
		// TODO: 添加模型
<%- end %>
	}
}

// GetAllIndexes 获取所有需要创建的复合索引
func GetAllIndexes() []db.IndexDefinition {
	return []db.IndexDefinition{
<%- if eq .GenerateExample "true" %>
		{
			Table: models.TableName<% .ExampleService %>,
			Name:  "idx_<% .ExampleService | snake %>s_tenant_code",
			// Columns 在 go-sqlbuilder v0.7 起是「列定义字符串」而非 []string：
			// 这样才能表达 "name DESC"、"(tenant_id, name)" 这类形式。
			// [EN] Since go-sqlbuilder v0.7 Columns is a column DEFINITION string rather than
			// a []string, so forms like "name DESC" or "(tenant_id, name)" can be expressed.
			// 编码在租户内唯一：协议中心的 <% .ExampleService %>CreateRequest 里
			// code 是必填且带 min/max 约束，业务语义上是租户内的业务主键
			// [EN] Code is unique per tenant: the hub's create request marks it required
			Columns: "(tenant_id, code)",
			Unique:  true,
		},
<%- else %>
		// TODO: 添加索引
<%- end %>
	}
}

// GetAllComments 获取所有表注释
func GetAllComments() []db.TableComment {
	return []db.TableComment{
<%- if eq .GenerateExample "true" %>
		{
			Table:   models.TableName<% .ExampleService %>,
			Comment: models.<% .ExampleService %>Model{}.TableComment(),
		},
<%- else %>
		// TODO: 添加表注释
<%- end %>
	}
}

// GetMigratorConfig 获取迁移器配置
func GetMigratorConfig() *db.MigratorConfig {
	return &db.MigratorConfig{
		Models:             GetAllModels(),
		Indexes:            GetAllIndexes(),
		Comments:           GetAllComments(),
		SkipIndexOnError:   true,
		SkipCommentOnError: true,
	}
}

// AutoMigrate 自动迁移所有数据表
func AutoMigrate(gormDB *gorm.DB) error {
	migrator := db.NewMigrator(gormDB, GetMigratorConfig())
	return migrator.AutoMigrate()
}

// CreateIndexes 创建额外的索引
func CreateIndexes(gormDB *gorm.DB) error {
	migrator := db.NewMigrator(gormDB, &db.MigratorConfig{
		Indexes:          GetAllIndexes(),
		SkipIndexOnError: true,
	})
	return migrator.CreateIndexes()
}

// AddTableComments 添加表注释
func AddTableComments(gormDB *gorm.DB) error {
	migrator := db.NewMigrator(gormDB, &db.MigratorConfig{
		Comments:           GetAllComments(),
		SkipCommentOnError: true,
	})
	return migrator.AddComments()
}
