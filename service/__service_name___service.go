/*
 * @Author: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @Date: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @LastEditors: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @LastEditTime: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @FilePath: \[[ .ProjectName ]]\service\[[ .ServiceName | snake ]]_service.go
 * @Description: [[ .ServiceName ]] 服务实现
 *
 * Copyright (c) [[ .Year ]] by [[ .AuthorName ]], All Rights Reserved.
 */
package service

import (
	"github.com/[[ .OrgName ]]/[[ .ProjectName ]]/repository"
)

// [[ .ServiceName ]]ServiceImpl 服务实现
type [[ .ServiceName ]]ServiceImpl struct {
	// Unimplemented[[ .ServiceName ]]ServiceServer
	repo *repository.RepositoryFactory
}

// New[[ .ServiceName ]]Service 创建服务实例
func New[[ .ServiceName ]]Service(repo *repository.RepositoryFactory) *[[ .ServiceName ]]ServiceImpl {
	return &[[ .ServiceName ]]ServiceImpl{
		repo: repo,
	}
}
