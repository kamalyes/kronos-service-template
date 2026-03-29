/*
 * @Author: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @Date: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @LastEditors: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @LastEditTime: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @FilePath: \[[ .ProjectName ]]\middleware\public_methods.go
 * @Description: 公共方法白名单
 *
 * Copyright (c) [[ .Year ]] by [[ .AuthorName ]], All Rights Reserved.
 */
package middleware

// publicMethods 完全公开的 gRPC 方法（无需任何身份信息）
var publicMethods = map[string]bool{
	// TODO: 添加公开方法
}

// IsPublicMethod 判断 gRPC FullMethod 是否完全公开
func IsPublicMethod(method string) bool {
	return publicMethods[method]
}
