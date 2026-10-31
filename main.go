/*
 * @Author: kamalyes example@qq.com
 * @Date: 2026-10-31 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-31 00:00:00
 * @FilePath: \kronos-my-service\main.go
 * @Description: MyService 服务启动器
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package main

import (
	"github.com/kamalyes/kronos-my-service/bootstrap"
	goconfig "github.com/kamalyes/go-config"
	"log"
	"os"
	"time"
)

func init() {
	goconfig.RegisterEnvPrefixes("laos", "laos")
	goconfig.RegisterEnvPrefixes("pakistan", "pakistan")
	goconfig.RegisterEnvPrefixes("myanmar", "myanmar")
}

func main() {
	time.Local = time.UTC

	server := bootstrap.NewMyServiceGateway()

	if err := server.Run(); err != nil {
		log.Fatalf("Gateway startup failed: %v", err)
		os.Exit(1)
	}
}
