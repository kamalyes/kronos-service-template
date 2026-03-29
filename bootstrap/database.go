/*
 * @Author: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @Date: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @LastEditors: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @LastEditTime: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @FilePath: \[[ .ProjectName ]]\bootstrap\database.go
 * @Description: 数据库初始化，执行自动迁移、索引创建和表注释
 *
 * Copyright (c) [[ .Year ]] by [[ .AuthorName ]], All Rights Reserved.
 */
package bootstrap

import (
	"fmt"
	"github.com/[[ .OrgName ]]/[[ .ProjectName ]]/migrations"
	goconfig "github.com/kamalyes/go-config"
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
)

// initDatabase 初始化数据库
func (s *[[ .ServiceName ]]Gateway) initDatabase() error {
	db := gwglobal.DB
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	gwglobal.LOGGER.Info("Running database migrations...")
	if err := migrations.AutoMigrate(db); err != nil {
		return fmt.Errorf("auto migrate failed: %v", err)
	}

	if err := migrations.CreateIndexes(db); err != nil {
		gwglobal.LOGGER.Warn(fmt.Sprintf("Index creation warning: %v", err))
	}

	if err := migrations.AddTableComments(db); err != nil {
		gwglobal.LOGGER.Warn(fmt.Sprintf("Table comments warning: %v", err))
	}

	env := goconfig.GetEnvironment()
	if env == goconfig.EnvLocal || env == goconfig.EnvDevelopment {
		if err := migrations.SeedData(db); err != nil {
			gwglobal.LOGGER.Warn(fmt.Sprintf("Seed data warning: %v", err))
		} else {
			gwglobal.LOGGER.Info("Seed data completed (dev environment)")
		}
	}

	gwglobal.LOGGER.Info("Database migrations completed")
	return nil
}
