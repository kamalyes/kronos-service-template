/*
 * @Author: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @Date: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @FilePath: \[[ .ProjectName ]]\main.go
 * @Description: [[ .ServiceName ]] 服务启动器
 *
 * Copyright (c) [[ .Year ]] by [[ .AuthorName ]], All Rights Reserved.
 */

package main

import (
	"github.com/[[ .OrgName ]]/[[ .ProjectName ]]/bootstrap"
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

	server := bootstrap.New[[ .ServiceName ]]Gateway()

	if err := server.Run(); err != nil {
		log.Fatalf("Gateway startup failed: %v", err)
		os.Exit(1)
	}
}
