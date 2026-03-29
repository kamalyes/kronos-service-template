/*
 * @Author: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @Date: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @LastEditors: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @LastEditTime: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @FilePath: \[[ .ProjectName ]]\migrations\auto_migrate.go
 * @Description: 数据库自动迁移 基于 go-sqlbuilder Migrator
 *
 * Copyright (c) [[ .Year ]] by [[ .AuthorName ]], All Rights Reserved.
 */

package migrations

import (
	"github.com/[[ .OrgName ]]/[[ .ProjectName ]]/models"
	"github.com/kamalyes/go-sqlbuilder/db"
	"gorm.io/gorm"
)

// GetAllModels 获取所有需要初始化的数据库模型
func GetAllModels() []interface{} {
	return []interface{}{
[[- if eq .GenerateExample "true" ]]
		&models.[[ .ExampleName ]]Model{},
[[- else ]]
		// TODO: 添加模型
[[- end ]]
	}
}

// GetAllIndexes 获取所有需要创建的复合索引
func GetAllIndexes() []db.IndexDefinition {
	return []db.IndexDefinition{
[[- if eq .GenerateExample "true" ]]
		{
			Table: models.TableName[[ .ExampleName ]],
			Name:  "idx_[[ .ExampleName | snake ]]s_name",
			Columns: []string{
				"name",
			},
			Unique: true,
		},
[[- else ]]
		// TODO: 添加索引
[[- end ]]
	}
}

// GetAllComments 获取所有表注释
func GetAllComments() []db.TableComment {
	return []db.TableComment{
[[- if eq .GenerateExample "true" ]]
		{
			Table:   models.TableName[[ .ExampleName ]],
			Comment: models.[[ .ExampleName ]]Model{}.TableComment(),
		},
[[- else ]]
		// TODO: 添加表注释
[[- end ]]
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
