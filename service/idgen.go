/*
 * @Author: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @Date: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @LastEditors: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @LastEditTime: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @FilePath: \[[ .ProjectName ]]\service\idgen.go
 * @Description: 服务层共享ID生成器（Snowflake）
 *
 * Copyright (c) [[ .Year ]] by [[ .AuthorName ]], All Rights Reserved.
 */
package service

import (
	"github.com/kamalyes/go-toolbox/pkg/idgen"
	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/kamalyes/go-toolbox/pkg/safe"
)

// snowflakeGen 进程内共享的 Snowflake 生成器
var snowflakeGen = idgen.NewSnowflakeGenerator(
	osx.GetWorkerIdForSnowflake(),
	osx.GetDatacenterId(),
)

// NewSnowflakeID 生成一个新的短ID字符串
func NewSnowflakeID() string {
	return safe.ShortHashWithLength(snowflakeGen.GenerateRequestID(), 8)
}
