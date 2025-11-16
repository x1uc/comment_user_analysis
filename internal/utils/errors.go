package utils

import (
	"errors"
	"fmt"
)

// AppError 应用错误类型
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 返回原始错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError 创建新的应用错误
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 预定义错误码
const (
	ErrCodeNetwork     = 1001
	ErrCodeParse       = 1002
	ErrCodeAuth        = 1003
	ErrCodeNotFound    = 1004
	ErrCodeRateLimit   = 1005
	ErrCodeConfig      = 1006
	ErrCodeExport      = 1007
)

// 预定义错误
var (
	ErrNoMoreData      = errors.New("没有更多数据")
	ErrInvalidResponse = errors.New("无效的响应")
	ErrUserNotFound    = errors.New("用户不存在")
)

// 便捷的构造函数
func NewNetworkError(message string, err error) *AppError {
	return NewAppError(ErrCodeNetwork, message, err)
}

func NewParseError(message string, err error) *AppError {
	return NewAppError(ErrCodeParse, message, err)
}

func NewAuthError(message string, err error) *AppError {
	return NewAppError(ErrCodeAuth, message, err)
}

func NewNotFoundError(message string, err error) *AppError {
	return NewAppError(ErrCodeNotFound, message, err)
}

func NewRateLimitError(message string, err error) *AppError {
	return NewAppError(ErrCodeRateLimit, message, err)
}

func NewConfigError(message string, err error) *AppError {
	return NewAppError(ErrCodeConfig, message, err)
}

func NewExportError(message string, err error) *AppError {
	return NewAppError(ErrCodeExport, message, err)
}