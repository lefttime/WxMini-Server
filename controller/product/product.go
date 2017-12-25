package product

import (
  "encoding/json"
  "gopkg.in/kataras/iris.v6"
  "github.com/lefttime/MyAssistant/model"
  "github.com/lefttime/MyAssistant/controller/common"
)

func FetchDiamondsInfo( ctx *iris.Context ) {
  SendErrJSON          := common.SendErrJSON
  SendSessionErrorJSON := common.SendSessionErrorJSON

  session    := ctx.Session()
  sessionKey := session.GetString( "wxAppSessionKey" )
  if sessionKey=="" {
    SendSessionErrorJSON( "Session错误", ctx )
    return
  }

  diamonds := model.FetchProductsByType( 0 )

  jsonData, err := json.Marshal( diamonds )
  if err != nil {
    SendErrJSON( "error", ctx )
    return
  }

  ctx.JSON( iris.StatusOK, iris.Map{
    "errNo" : model.ErrorCode.Success,
    "msg"   : "success",
    "data"  : string(jsonData),
  })
}