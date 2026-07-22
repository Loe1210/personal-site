package xhttp

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/Loe1210/personal-site/internal/xerrors"
)

type Envelope struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func OK(c *app.RequestContext, data any) {
	c.JSON(consts.StatusOK, Envelope{Code: xerrors.CodeOK, Msg: "success", Data: data})
}

func Fail(c *app.RequestContext, err error) {
	c.JSON(consts.StatusOK, Envelope{Code: xerrors.CodeOf(err), Msg: xerrors.MessageOf(err), Data: nil})
}
