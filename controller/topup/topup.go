package topup

import (
  "net"
  "time"
  "strings"
  "encoding/json"
  "gopkg.in/kataras/iris.v6"
  "github.com/lefttime/MyAssistant/model"
  "github.com/lefttime/MyAssistant/controller/game"
  "github.com/lefttime/MyAssistant/controller/common"
)

func FetchTopupRecentInfo( ctx *iris.Context ) {
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

  type RecentInfo struct {
    TodayDiamonds    int                        `json:"today"`
    LastWeekDiamonds int                        `json:"lastweek"`
    Presentees       []model.TopupPresenteeInfo `json:"presentees"`
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

  var recentInfo RecentInfo
  recentInfo.TodayDiamonds    = countTopupsDiamonds( model.FetchUserTodayTopups( userId )    )
  recentInfo.LastWeekDiamonds = countTopupsDiamonds( model.FetchUserLastWeekTopups( userId ) )
  recentInfo.Presentees       = queryPresenteeTopups( userId, 0 )

  jsonData, err := json.Marshal( recentInfo )
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

func TopupDiamondForUser( ctx *iris.Context ) {
  SendErrJSON          := common.SendErrJSON
  SendSessionErrorJSON := common.SendSessionErrorJSON

  session    := ctx.Session()
  sessionKey := session.GetString( "wxAppSessionKey" )
  if sessionKey=="" {
    SendSessionErrorJSON( "Session错误", ctx )
    return
  }

  type RequestData struct {
    PresenteeUnionID string `json:"presenteeUnionId"`
    PresenteeOpenID  string `json:"presenteeOpenId"`
    Diamond          uint   `json:"diamond"`
    WithdrawToGame   uint   `json:"withdrawToGame"`
  }

  var wxRequestData RequestData
  userOpenId := session.GetString( "wxAppOpenID" )

  if ctx.ReadJSON( &wxRequestData ) != nil || userOpenId=="" {
    SendErrJSON( "参数错误", ctx )
    return
  }

  userInfo := model.FetchUserInfo( userOpenId )
  if userInfo.Diamond < int( wxRequestData.Diamond ) {
    SendErrJSON( "钻石不足", ctx )
    return
  }

  onlyGamePlayer := false
  presenteeInfo := model.FetchUserDetail( wxRequestData.PresenteeOpenID )
  if presenteeInfo.ID==0 {
    onlyGamePlayer = true
    presenteeInfo = game.SearchUserInfoByUnionId( wxRequestData.PresenteeUnionID )
    if presenteeInfo.ID==0 {
      SendErrJSON( "玩家不存在", ctx )
      return
    }
  }

  var result model.Topup
  result.UserID      = userInfo.ID
  result.PresenteeID = presenteeInfo.ID
  result.Count       = wxRequestData.Diamond
  result.CreatedAt   = time.Now()
  result.Status      = 1

  // 未完善，待后期作异步处理
  if wxRequestData.WithdrawToGame==1 {
    if !withdrawDiamondToGame( int(result.ID), int(wxRequestData.Diamond), presenteeInfo.UnionID, userInfo.Nickname ) {
      SendErrJSON( "充值失败", ctx )
      return
    }
  } else {
    if !onlyGamePlayer {
      model.IncreaseDiamond( wxRequestData.PresenteeOpenID, int(wxRequestData.Diamond) )  
    }
  }
  model.DB.Create( &result )
  model.DecreaseDiamond( userOpenId, int(wxRequestData.Diamond) )

  ctx.JSON( iris.StatusOK, iris.Map{
    "errNo" : model.ErrorCode.Success,
    "msg"   : "success",
    "data"  : iris.Map{},
  })
}

func withdrawDiamondToGame( reqId int, diamond int, openId string, operator string ) bool {
  type RequestData struct {
    ReqId    int    `json:"ReqId"`
    MsgId    int    `json:"MsgId"`
    OpenId   string `json:"OpenId"`
    CardType int    `json:"CardType"`
    CardNum  int    `json:"CardNum"`
    OperType int    `json:"OperType"`
    GameType int    `json:"GameType"`
    OperName string `json:"OperName"`
  }

  var requestData RequestData
  requestData.ReqId    = reqId
  requestData.MsgId    = 0
  requestData.OpenId   = openId
  requestData.CardType = 2
  requestData.CardNum  = diamond
  requestData.OperType = 1
  requestData.GameType = 10003
  requestData.OperName = operator

  jsonData, err := json.Marshal( requestData )
  if err != nil {
    return false
  }

  conn, err := net.DialTimeout( "tcp", "datum.moy2017.com:36336", 2 * time.Second )
  if err != nil {
    return false
  }

  conn.Write( jsonData )
  defer conn.Close()
  var buf = make( []byte, 128 )
  total, err := conn.Read( buf )
  if err != nil {
    return false
  }
  buf = buf[:total]

  type MsgInfo struct {
    ErrorCode int    `json:"errorCode"`
    ErrorMsg  string `json:"errorMsg"`
    ReqId     int    `json:"ReqId"`
  }

  bufStr := string( buf )
  bufStr = strings.Replace( bufStr, "'", "\"", -1 )

  var msgInfo MsgInfo
  if err := json.Unmarshal( []byte( bufStr ), &msgInfo ); err != nil {
    return false
  }

  return true
}

func countTopupsDiamonds( topups []model.Topup ) int {
  var result int
  for _, value := range topups {
    result = result + int(value.Count)
  }

  return result
}

func queryPresenteeTopups( userId int, limit int ) []model.TopupPresenteeInfo {
  topups := model.FetchTopupsByUserId( userId, limit )

  result := make( []model.TopupPresenteeInfo, len( topups) )
  for idx := 0; idx < len( topups ); idx++ {
    user := model.FetchUserDetailById( int( topups[idx].PresenteeID ) )
    result[idx].OpenID    = user.OpenID
    result[idx].Avatar    = user.Avatar
    result[idx].Nickname  = user.Nickname
    result[idx].Post      = user.Post
    result[idx].GameID    = user.GameID
    result[idx].Phone     = user.Phone
    result[idx].Count     = topups[idx].Count
    result[idx].CreatedAt = model.FormatDatetime( topups[idx].CreatedAt, false )
  }

  return result
}
