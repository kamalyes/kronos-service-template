/*
 * @Author: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @Date: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @LastEditors: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @LastEditTime: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @FilePath: \[[ .ProjectName ]]\errors\grpc.go
 * @Description: gRPC 错误转换处理
 *
 * Copyright (c) [[ .Year ]] by [[ .AuthorName ]], All Rights Reserved.
 */

package errors

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

// ToGRPCError 将业务错误代码转换为标准 gRPC 错误
func ToGRPCError(ctx context.Context, bizErrorCode string) error {
	errCode := NewI18nError(ctx, bizErrorCode)
	grpcCode := MapBizCodeToGRPCCode(bizErrorCode)
	return status.Error(grpcCode, errCode.Message)
}

// ToGRPCErrorWithTemplate 将带模板的业务错误代码转换为标准 gRPC 错误
func ToGRPCErrorWithTemplate(ctx context.Context, bizErrorCode string, templateData map[string]interface{}) error {
	errCode := NewI18nErrorWithTemplate(ctx, bizErrorCode, templateData)
	grpcCode := MapBizCodeToGRPCCode(bizErrorCode)
	return status.Error(grpcCode, errCode.Message)
}

// MapBizCodeToGRPCCode 将业务错误代码映射到 gRPC 标准错误码
func MapBizCodeToGRPCCode(bizErrorCode string) codes.Code {
	switch bizErrorCode {
	case BizErrCodeInvalidArgument, BizErrCodeValidationFailed, BizErrCodeInvalidParameter:
		return codes.InvalidArgument
	case BizErrCodeUnauthenticated:
		return codes.Unauthenticated
	case BizErrCodePermissionDenied:
		return codes.PermissionDenied
	case BizErrCodeConflict:
		return codes.AlreadyExists
	}

	if bizErrorCode == BizErrCodeNotFound || strings.HasSuffix(bizErrorCode, "_not_found") {
		return codes.NotFound
	}

	return codes.Internal
}
