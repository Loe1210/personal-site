package response

import "github.com/Loe1210/personal-site/pkg/errno"

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) Body {
	return Body{
		Code:    errno.Success.Code,
		Message: errno.Success.Message,
		Data:    data,
	}
}

func Error(code int, message string) Body {
	return Body{
		Code:    code,
		Message: message,
	}
}

func AppError(err *errno.AppError) Body {
	return Body{
		Code:    err.Code,
		Message: err.Message,
	}
}

func ErrorWithMessage(err *errno.AppError, message string) Body {
	return Body{
		Code:    err.Code,
		Message: message,
	}
}
