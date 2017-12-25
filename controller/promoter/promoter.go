package promoter

import (
  "encoding/json"
  "gopkg.in/kataras/iris.v6"
  "github.com/lefttime/MyAssistant/model"
  "github.com/lefttime/MyAssistant/controller/common"
)

func FetchPromoterInfo( ctx *iris.Context ) {
  SendErrJSON          := common.SendErrJSON
  SendSessionErrorJSON := common.SendSessionErrorJSON

  session    := ctx.Session()
  sessionKey := session.GetString( "wxAppSessionKey" )
  if sessionKey=="" {
    SendSessionErrorJSON( "Session错误", ctx )
    return
  }

  var openId = session.GetString( "wxAppOpenID" )
  if openId=="" {
    SendErrJSON( "参数错误, 未找到OpenID", ctx )
    return
  }

  userId := model.GetUserIdByOpenId( openId )
  if userId <= 0 {
    ctx.JSON( iris.StatusOK, iris.Map{
      "errNo" : model.ErrorCode.NotFound,
      "msg"   : "failed",
      "data"  : iris.Map{},
    })
    return
  }

  var promoterInfo = queryPromoterInfo( userId, 0 )
  jsonData, err := json.Marshal( promoterInfo )
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

func queryPromoterInfo( userId int, limit int ) []model.PromoterInfo {
  promoters := model.FetchPromotersByLeaderID( userId, limit )

  result := make( []model.PromoterInfo, len( promoters) )
  for idx := 0; idx < len( promoters ); idx++ {
    user := model.FetchUserDetailById( int( promoters[idx].UserID ) )
    result[idx].OpenID   = user.OpenID
    result[idx].Avatar   = user.Avatar
    result[idx].Nickname = user.Nickname
    result[idx].Diamond  = user.Diamond
    result[idx].Phone    = user.Phone
  }

  return result
}