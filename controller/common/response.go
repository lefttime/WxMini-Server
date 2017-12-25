package common

import (
  "gopkg.in/kataras/iris.v6"
  "github.com/lefttime/MyAssistant/model"
)

func SendErrJSON( msg string, ctx *iris.Context ) {
  ctx.JSON( iris.StatusOK, iris.Map{
    "errNo" : model.ErrorCode.Error,
    "msg"   : msg,
    "data"  : iris.Map{},
  })
}

func SendSessionErrorJSON( msg string, ctx *iris.Context ) {
  ctx.JSON( iris.StatusOK, iris.Map{
    "errNo" : model.ErrorCode.SessionError,
    "msg"   : msg,
    "data"  : iris.Map{},
  })
}