/*
 * @Author: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @Date: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @LastEditors: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @LastEditTime: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @FilePath: \[[ .ProjectName ]]\errors\errors.go
 * @Description: 国际化错误处理核心函数
 *
 * Copyright (c) [[ .Year ]] by [[ .AuthorName ]], All Rights Reserved.
 */

package errors

import (
	"context"
	"errors"
	commonpb "[[ .ShareProtoModule ]]/pb/common"
	"github.com/kamalyes/go-rpc-gateway/middleware"
)

// New 创建一个新的错误
func New(message string) error {
	return errors.New(message)
}

// NewI18nError 创建国际化错误响应
func NewI18nError(ctx context.Context, errorCode string) *commonpb.Error {
	message := middleware.GetMsgByKey(ctx, errorCode)
	if message == "" || message == errorCode {
		message = errorCode
	}
	return &commonpb.Error{
		Code:    errorCode,
		Message: message,
	}
}

// NewI18nErrorWithTemplate 创建带模板数据的国际化错误响应
func NewI18nErrorWithTemplate(ctx context.Context, errorCode string, templateData map[string]interface{}) *commonpb.Error {
	message := middleware.GetMsgWithMap(ctx, errorCode, templateData)
	if message == "" || message == errorCode {
		message = errorCode
	}
	return &commonpb.Error{
		Code:    errorCode,
		Message: message,
	}
}
