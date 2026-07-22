package xerrors

import "errors"

const (
	CodeOK                     int32 = 0
	CodeInvalidArgument        int32 = 20010001
	CodeInternal               int32 = 20010002
	CodeAuthLoginRequired      int32 = 20020001
	CodeAuthSessionExpired     int32 = 20020002
	CodeAuthPermissionDenied   int32 = 20020003
	CodeAuthUpstreamFailed     int32 = 20020004
	CodeContentArticleNotFound int32 = 20030001
)

type AppError struct {
	code int32
	msg  string
}

func New(code int32, msg string) *AppError {
	return &AppError{code: code, msg: msg}
}

func (e *AppError) Error() string {
	return e.msg
}

func (e *AppError) Code() int32 {
	return e.code
}

func (e *AppError) Message() string {
	return e.msg
}

func CodeOf(err error) int32 {
	if err == nil {
		return CodeOK
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code()
	}
	return CodeInternal
}

func MessageOf(err error) string {
	if err == nil {
		return "success"
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Message()
	}
	return "internal error"
}
