/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02" %> 17:05:08
 * @LastEditors: <% .AuthorName %> <% .AuthorEmail %>
 * @LastEditTime: <% nowFmt "2006-01-02" %> 19:28:30
 * @FilePath: \<% .ProjectName %>\migrations\migrate.go
 * @Description: 数据库种子数据迁移
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */

package migrations

import (
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
	"gorm.io/gorm"
)

// SeedData 初始化种子数据
func SeedData(db *gorm.DB) error {
	gwglobal.LOGGER.Info("Seeding initial data...")

	// TODO: 添加种子数据初始化逻辑

	gwglobal.LOGGER.Info("Seed data completed")
	return nil
}
