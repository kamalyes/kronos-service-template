/*
 * @Author: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @Date: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @LastEditors: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @LastEditTime: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @FilePath: \[[ .ProjectName ]]\bootstrap\subscriber.go
 * @Description: 订阅者管理，处理分布式事件订阅
 *
 * Copyright (c) [[ .Year ]] by [[ .AuthorName ]], All Rights Reserved.
 */
package bootstrap

import (
	gwglobal "github.com/kamalyes/go-rpc-gateway/global"
)

// startSubscribers 启动所有事件订阅者
func (s *[[ .ServiceName ]]Gateway) startSubscribers() {
	// TODO: 添加事件订阅逻辑
	gwglobal.LOGGER.Info("All subscribers started")
}
