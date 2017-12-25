package user

import (
  "fmt"
  "time"
  "net/http"
  "strings"
  "strconv"
  "encoding/json"
  "gopkg.in/kataras/iris.v6"
  "github.com/lefttime/MyAssistant/model"
  "github.com/lefttime/MyAssistant/utils"
  "github.com/lefttime/MyAssistant/config"
  "github.com/lefttime/MyAssistant/controller/game"
  "github.com/lefttime/MyAssistant/controller/order"
  "github.com/lefttime/MyAssistant/controller/common"
)

func WxAppLogin( ctx *iris.Context ) {
  SendErrJSON := common.SendErrJSON
  code := ctx.FormValue( "code" )
  if code=="" {
    SendErrJSON( "code不能为空", ctx )
    return
  }

  CodeToSessURL := config.WxAppConfig.CodeToSessURL
  CodeToSessURL  = strings.Replace( CodeToSessURL, "{code}", code, -1 )

  resp, err := http.Get( CodeToSessURL )
  if err != nil {
    SendErrJSON( "请求微信授权出错", ctx )
    return
  }
  defer resp.Body.Close()

  if resp.StatusCode != 200 {
    SendErrJSON( "请求微信授权出错", ctx )
    return
  }

  var data map[string]interface{}
  err = json.NewDecoder( resp.Body ).Decode( &data )
  if err != nil {
    SendErrJSON( "微信授权数据解密出错", ctx )
    return
  }

  if _, ok := data["session_key"]; !ok {
    SendErrJSON( "微信授权数据不完整", ctx )
    return
  }

  var openID     string
  var sessionKey string
  openID     = data["openid"].( string )
  sessionKey = data["session_key"].( string )
  session   := ctx.Session()
  session.Set( "wxAppOpenID",     openID     )
  session.Set( "wxAppSessionKey", sessionKey )

  resData := iris.Map{}
  resData[config.ServerConfig.SessionID] = session.ID()
  ctx.JSON( iris.StatusOK, iris.Map{
    "errNo" : model.ErrorCode.Success,
    "msg"   : "success",
    "data"  : resData,
  })
}

func SetWxAppUserInfo( ctx *iris.Context ) {
  SendErrJSON          := common.SendErrJSON
  SendSessionErrorJSON := common.SendSessionErrorJSON

  type EncryptedUser struct {
    EncryptedData string `json:"encryptedData"`
    IV            string `json:"iv"`
    Platform      string `json:"platform"`
  }

  var wxAppUser EncryptedUser

  if ctx.ReadJSON( &wxAppUser ) != nil {
    SendErrJSON( "参数错误", ctx )
    return
  }
  session    := ctx.Session()
  sessionKey := session.GetString( "wxAppSessionKey" )
  if sessionKey=="" {
    SendSessionErrorJSON( "Session错误", ctx )
    return
  }

  userInfoStr, err := utils.DecodeWxAppUserInfo( wxAppUser.EncryptedData, sessionKey, wxAppUser.IV )
  if err != nil {
    fmt.Println( err.Error() )
    SendErrJSON( "error", ctx )
    return
  }

  var wxUser model.WxAppUser
  if err := json.Unmarshal( []byte( userInfoStr ), &wxUser ); err != nil {
    SendErrJSON( "error", ctx )
    return
  }

  session.Set( "wxAppUser", wxUser )
  ctx.JSON( iris.StatusOK, iris.Map{
    "errNo" : model.ErrorCode.Success,
    "msg"   : "success",
    "data"  : iris.Map{},
  })

  var user model.User
  user.CreatedAt = time.Now()
  user.UpdatedAt = user.CreatedAt
  user.OpenID    = wxUser.OpenID
  user.UnionID   = wxUser.UnionID
  user.Nickname  = wxUser.Nickname
  user.Gender    = wxUser.Gender
  user.Avatar    = wxUser.AvatarURL
  user.City      = wxUser.City
  user.Province  = wxUser.Province
  user.Country   = wxUser.Country
  user.Platform  = wxAppUser.Platform
  model.UpdateUserInfo( user )

  return
}

func FetchUserInfo( ctx *iris.Context ) {
  SendErrJSON          := common.SendErrJSON
  SendSessionErrorJSON := common.SendSessionErrorJSON

  type RequestData struct {
    OpenID  string `json:"openId"`
    UnionID string `json"unionId"`
  }

  var wxFetchWxAppUser RequestData

  if ctx.ReadJSON( &wxFetchWxAppUser ) != nil {
    SendErrJSON( "参数错误", ctx )
    return
  }

  session    := ctx.Session()
  sessionKey := session.GetString( "wxAppSessionKey" )
  if sessionKey=="" {
    SendSessionErrorJSON( "Session错误", ctx )
    return
  }

  var userInfo model.UserDetail
  if wxFetchWxAppUser.UnionID != "" {
    userInfo = model.FetchUserDetailByUnionId( wxFetchWxAppUser.UnionID )
    if userInfo.ID==0 {
      userInfo = game.SearchUserInfoByUnionId( wxFetchWxAppUser.UnionID )
    }
  } else {
    if wxFetchWxAppUser.OpenID=="" {
      wxFetchWxAppUser.OpenID = session.GetString( "wxAppOpenID" )
    }
    userInfo = model.FetchUserDetail( wxFetchWxAppUser.OpenID )
  }

  if userInfo.ID != 0 {
    jsonData, err := json.Marshal( userInfo )
    if err != nil {
      SendErrJSON( "数据封装错误", ctx )
    } else {
      ctx.JSON( iris.StatusOK, iris.Map {
        "errNo" : model.ErrorCode.Success,
        "msg"   : "success",
        "data"  : string( jsonData ),
      } ) 
    }
  } else {
    ctx.JSON( iris.StatusOK, iris.Map {
      "errNo" : model.ErrorCode.NotFound,
      "msg"   : "failed",
      "data"  : iris.Map{},
    } )
  }
}

func Payment( ctx *iris.Context ) {
  SendErrJSON          := common.SendErrJSON
  SendSessionErrorJSON := common.SendSessionErrorJSON

  session    := ctx.Session()
  sessionKey := session.GetString( "wxAppSessionKey" )
  if sessionKey=="" {
    SendSessionErrorJSON( "Session错误", ctx )
    return
  }

  type RequestData struct {
    OpenID       string `json:"openId"`
    ProductID    uint   `json:"productId"`
    ProductCount uint   `json:"productCount"`
    TotalPrice   uint   `json:"totalPrice"`
  }

  var wxPaymentData RequestData
  wxPaymentData.OpenID = session.GetString( "wxAppOpenID" )

  if ctx.ReadJSON( &wxPaymentData ) != nil || wxPaymentData.OpenID=="" {
    SendErrJSON( "参数错误", ctx )
    return
  }

  userInfo := model.FetchUserInfo( wxPaymentData.OpenID )
  order.GenerateOrder( userInfo.ID,
                       userInfo.Post,
                       userInfo.Level,
                       wxPaymentData.ProductID,
                       wxPaymentData.ProductCount,
                       wxPaymentData.TotalPrice )
  userInfo.Diamond = userInfo.Diamond + int( wxPaymentData.ProductCount )
  userInfo.Exp     = userInfo.Exp + wxPaymentData.ProductCount / 10
  userInfo.Level   = model.CalcLevelByExperience( userInfo.Exp, userInfo.Post )
  model.ModifyUserInfo( userInfo )

  // 仅供测试用

  jsonData, err := json.Marshal( userInfo )
  if err != nil {
    SendErrJSON( "数据封装错误", ctx )
  } else {
    ctx.JSON( iris.StatusOK, iris.Map{
      "errNo" : model.ErrorCode.Success,
      "msg"   : "success",
      "data"  : string( jsonData ),
    }) 
  }
}

func SearchUserInfo( ctx *iris.Context ) {
  SendErrJSON          := common.SendErrJSON
  SendSessionErrorJSON := common.SendSessionErrorJSON
  
  type RequestData struct {
    ID     string `json:"id"`
    OpenID string `json:"openId"`
  }

  var wxRequestData RequestData

  if ctx.ReadJSON( &wxRequestData ) != nil {
    SendErrJSON( "参数错误", ctx )
    return
  }

  idLen := len( wxRequestData.ID )
  value_int, err := strconv.Atoi( wxRequestData.ID )
  if idLen==0 || idLen > 11 || err != nil {
    SendErrJSON( "搜索ID错误", ctx )
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

  if !checkPermission( wxRequestData.OpenID ) {
    SendErrJSON( "无权限查找用户信息", ctx )
    return
  }

  var userInfo model.UserDetail
  if idLen==11 {
    userInfo = model.FetchUserDetailByPhone( wxRequestData.ID )
  } else {
    userInfo = game.SearchUserInfo( value_int )
  }

  if userInfo.ID==0 {
    SendErrJSON( "查找不到用户信息", ctx )
  } else {
    jsonData, err := json.Marshal( userInfo )
    if err != nil {
      SendErrJSON( "查找不到用户信息", ctx )
    } else {
      ctx.JSON( iris.StatusOK, iris.Map{
        "errNo" : model.ErrorCode.Success,
        "msg"   : "success",
        "data"  : string( jsonData ),
      }) 
    }
  }
}

func checkPermission( openId string ) bool {
  userInfo := model.FetchUserDetail( openId )
  return userInfo.ID != 0 && userInfo.Post > 1
}