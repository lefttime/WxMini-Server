package game

import (
  "time"
  "github.com/lefttime/MyAssistant/model"
)

type UserGameInfo struct {
  ID        uint   `json:"id"`
  Nickname  string `json:"nickname"`
  Gender    int    `json:"gender"`
  Province  string `json:"province"`
  City      string `json:"city"`
  Country   string `json:"country"`
  Avatar    string `json:"avatar"`
  UnionID   string `json:"unionId"`
  Diamond   int    `json:"diamond"`
  CreatedAt int    `json:"createdAt"`
  UpdatedAt int    `json:"updatedAt"`
}

func SearchUserInfo( gameId int ) model.UserDetail {
  var userInfo model.User

  var gameInfo UserGameInfo
  model.GameDB.Raw( "select id, Nike as nickname, Sex as gender, Provice as province, City as city, Country as country, HeadImageUrl as avatar, UnionId as union_id, NumsCard2 as diamond, RegTime as created_at, LastLoginTime as updated_at from user where id = ?", gameId ).Scan( &gameInfo )
  userInfo.ID        = gameInfo.ID
  userInfo.Nickname  = gameInfo.Nickname
  userInfo.Gender    = gameInfo.Gender
  userInfo.Province  = gameInfo.Province
  userInfo.City      = gameInfo.City
  userInfo.Country   = gameInfo.Country
  userInfo.Avatar    = gameInfo.Avatar
  userInfo.UnionID   = gameInfo.UnionID
  userInfo.Diamond   = gameInfo.Diamond
  userInfo.CreatedAt = time.Unix( int64(gameInfo.CreatedAt), 0 )
  userInfo.UpdatedAt = time.Unix( int64(gameInfo.UpdatedAt), 0 )

  return model.UserDetail{ userInfo, model.FormatDatetime( userInfo.CreatedAt, false ), model.FormatDatetime( userInfo.UpdatedAt, false ), gameInfo.ID }
}

func SearchUserInfoByUnionId( unionId string ) model.UserDetail {
  var userInfo model.User

  var gameInfo UserGameInfo
  model.GameDB.Raw( "select id, Nike as nickname, Sex as gender, Provice as province, City as city, Country as country, HeadImageUrl as avatar, UnionId as union_id, NumsCard2 as diamond, RegTime as created_at, LastLoginTime as updated_at from user where unionid = ?", unionId ).Scan( &gameInfo )
  userInfo.ID        = gameInfo.ID
  userInfo.Nickname  = gameInfo.Nickname
  userInfo.Gender    = gameInfo.Gender
  userInfo.Province  = gameInfo.Province
  userInfo.City      = gameInfo.City
  userInfo.Country   = gameInfo.Country
  userInfo.Avatar    = gameInfo.Avatar
  userInfo.UnionID   = gameInfo.UnionID
  userInfo.Diamond   = gameInfo.Diamond
  userInfo.CreatedAt = time.Unix( int64(gameInfo.CreatedAt), 0 )
  userInfo.UpdatedAt = time.Unix( int64(gameInfo.UpdatedAt), 0 )

  return model.UserDetail{ userInfo, model.FormatDatetime( userInfo.CreatedAt, false ), model.FormatDatetime( userInfo.UpdatedAt, false ), gameInfo.ID }
}