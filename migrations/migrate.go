/*
 * @Author: kamalyes example@qq.com
 * @Date: 2026-10-31 00:00:00
 * @LastEditors: kamalyes example@qq.com
 * @LastEditTime: 2026-10-31 00:00:00
 * @FilePath: \kronos-my-service\migrations\migrate.go
 * @Description: 数据库种子数据迁移
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
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
