/*
 * @Author: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @Date: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @LastEditors: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @LastEditTime: [[ nowFmt "2006-01-02" ]] 00:00:00
 * @FilePath: \[[ .ProjectName ]]\errors\codes.go
 * @Description: 业务错误代码常量定义
 *
 * Copyright (c) [[ .Year ]] by [[ .AuthorName ]], All Rights Reserved.
 */

package errors

const (
	// 通用错误
	BizErrCodeInternalError    = "error.internal"
	BizErrCodeInvalidArgument  = "error.invalid_argument"
	BizErrCodeUnauthenticated  = "error.unauthenticated"
	BizErrCodePermissionDenied = "error.permission_denied"
	BizErrCodeNotFound         = "error.not_found"
	BizErrCodeConflict         = "error.conflict"
	BizErrCodeValidationFailed = "error.validation_failed"
	BizErrCodeInvalidParameter = "error.invalid_parameter"
	BizErrCodeDatabaseOperation = "error.database_operation_failed"
	BizErrCodeConversionFailed  = "error.conversion_failed"
)
