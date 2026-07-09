package response

import (
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/Loe1210/personal-site/pkg/errno"
)

func WriteError(c *app.RequestContext, err *errno.AppError) {
	c.JSON(err.HTTPStatus, AppError(err))
}

func WriteErrorMessage(c *app.RequestContext, err *errno.AppError, message string) {
	c.JSON(err.HTTPStatus, ErrorWithMessage(err, message))
}

func WriteSuccess(c *app.RequestContext, data interface{}) {
	c.JSON(200, Success(data))
}
