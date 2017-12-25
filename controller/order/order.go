package order

import (
  "time"
  "encoding/json"
  "gopkg.in/kataras/iris.v6"
  "github.com/lefttime/MyAssistant/model"
  "github.com/lefttime/MyAssistant/controller/rebate"
  "github.com/lefttime/MyAssistant/controller/common"
)

func FetchOrdersInfo( ctx *iris.Context ) {
  SendErrJSON          := common.SendErrJSON
  SendSessionErrorJSON := common.SendSessionErrorJSON

  type RequestData struct {
    OpenID string `json:"openId"`
  }

  var wxRequestData RequestData
  if ctx.ReadJSON( &wxRequestData ) != nil {
    SendErrJSON( "参数错误", ctx )
    return
  }

  session    := ctx.Session()
  sessionKey := session.GetString( "wxAppSessionKey" )
  if sessionKey=="" {
    SendSessionErrorJSON( "Session错误", ctx )
    return
  }

  if wxRequestData.OpenID=="" {
    wxRequestData.OpenID = session.GetString( "wxAppOpenID" )
  }

  userId := model.GetUserIdByOpenId( wxRequestData.OpenID )
  if userId <= 0 {
    ctx.JSON( iris.StatusOK, iris.Map{
      "errNo" : model.ErrorCode.NotFound,
      "msg"   : "failed",
      "data"  : iris.Map{},
    })
    return
  }

  orderInfo := model.FetchOrdersByUserId( userId )
  jsonData, err := json.Marshal( orderInfo )
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

func GenerateOrder( userId uint, post uint, level uint, productId uint, productCount uint, totalPrice uint ) model.Order {
  var order model.Order
  order.UserID       = userId
  order.ProductID    = productId
  order.ProductCount = productCount
  order.Rebate       = rebate.CalcRebate( post, level, productCount )
  order.Discount     = 0          // 测试用
  order.OriginPrice  = float32( totalPrice )
  order.TotalPrice   = float32( totalPrice )
  order.Status       = 0          // 测试用
  order.CreatedAt    = time.Now()
  order.PayAt        = time.Now() // 测试用
  model.DB.Create( &order )

  return order
}
