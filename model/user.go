package model

import (
	"time"
)

type WxAppUser struct {
	OpenID    string `json:"openId"`
	UnionID   string `json:"unionId"`
	Nickname  string `json:"nickname"`
	Gender    int    `json:"gender"`
	AvatarURL string `json:"avatarUrl"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
}

type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	OpenID    string    `json:"openId"`
	UnionID   string    `json:"unionId"`
	Nickname  string    `json:"nickname"`
	Phone     string    `json:"phone"`
	Token     string    `json:"token"`
	Avatar    string    `json:"avatarUrl"`
	Gender    int       `json:"gender"`
	Post      uint      `json:"post"`
	Platform  string    `json:"platform"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	City      string    `json:"city"`
	Province  string    `json:"province"`
	Country   string    `json:"country"`
	Diamond   int       `json:"diamond"`
  Exp       uint      `json:"experience"`
  Level     uint      `json:"level"`
}

type UserDetail struct {
	User
  CreatedAt string `json:"createdAt"`
  UpdatedAt string `json:"updatedAt"`
  GameID   uint    `json:"gameId"`
}

func UpdateUserInfo( user User ) {
	var count int
	DB.Model( &User{} ).Where( "open_id = ?", user.OpenID ).Count( &count )
	if count==0 {
		user.CreatedAt = time.Now()
		user.UpdatedAt = user.CreatedAt
		DB.Create( &user )
	} else {
		DB.Model( &user ).Updates( map[string]interface{}{ "token": user.Token, "platform": user.Platform } )
	}
}

func ModifyUserInfo( user User ) {
  DB.Model( &user ).Updates( map[string]interface{}{ "nickname": user.Nickname,
                                                     "phone":    user.Phone,
                                                     "avatar":   user.Avatar,
                                                     "post":     user.Post,
                                                     "platform": user.Platform,
                                                     "city":     user.City,
                                                     "province": user.Province,
                                                     "country":  user.Country,
                                                     "diamond":  user.Diamond,
                                                     "exp":      user.Exp,
                                                     "level":    user.Level } )
}

func FetchUserInfo( openId string ) User {
	var result User
  DB.Where( "open_id = ?", openId ).First( &result )
  return result
}

func FetchUserInfoById( id int ) User {
	var result User
	DB.Where( "id = ?", id ).First( &result )

	return result
}

func FetchUserInfoByPhone( phone string ) User {
  var result User
  DB.Where( "phone = ?", phone ).First( &result )

  return result
}

func FetchUserInfoByUnionId( unionId string ) User {
  var result User
  DB.Where( "union_id = ?", unionId ).First( &result )

  return result
}

func FetchUserDetail( openId string ) UserDetail {
	var result UserDetail

  userInfo := FetchUserInfo( openId )
  GameDB.Raw( "select id as game_id from user where unionid = ?", userInfo.UnionID ).Scan( &result )

  return UserDetail{ userInfo, FormatDatetime( userInfo.CreatedAt, false ), FormatDatetime( userInfo.UpdatedAt, false ), result.GameID }
}

func FetchUserDetailById( id int ) UserDetail {
	var result UserDetail

  userInfo := FetchUserInfoById( id )
  GameDB.Raw( "select id as game_id from user where unionid = ?", userInfo.UnionID ).Scan( &result )

  return UserDetail{ userInfo, FormatDatetime( userInfo.CreatedAt, false ), FormatDatetime( userInfo.UpdatedAt, false ), result.GameID }
}

func FetchUserDetailByPhone( phone string ) UserDetail {
  var result UserDetail

  userInfo := FetchUserInfoByPhone( phone )
  GameDB.Raw( "select id as game_id from user where unionid = ?", userInfo.UnionID ).Scan( &result )

  return UserDetail{ userInfo, FormatDatetime( userInfo.CreatedAt, false ), FormatDatetime( userInfo.UpdatedAt, false ), result.GameID }
}

func FetchUserDetailByUnionId( unionId string ) UserDetail {
  var result UserDetail

  userInfo := FetchUserInfoByUnionId( unionId )
  GameDB.Raw( "select id as game_id from user where unionid = ?", userInfo.UnionID ).Scan( &result )

  return UserDetail{ userInfo, FormatDatetime( userInfo.CreatedAt, false ), FormatDatetime( userInfo.UpdatedAt, false ), result.GameID }
}

func GetUserIdByOpenId( openId string ) int {
	var user User
	DB.Where( "open_id = ?", openId ).First( &user )
	if user.ID==0 {
		return -1;
	}

	return int( user.ID )
}

func IncreaseDiamond( openId string, diamond int ) {
  var user User
  DB.Where( "open_id = ?", openId ).First( &user )
  if user.ID != 0 {
    DB.Model( &user ).Updates( map[string]interface{}{ "diamond": user.Diamond + diamond } )
  }
}

func DecreaseDiamond( openId string, diamond int ) {
  var user User
  DB.Where( "open_id = ?", openId ).First( &user )
  if user.ID != 0 {
    DB.Model( &user ).Updates( map[string]interface{}{ "diamond": user.Diamond - diamond } )
  }
}

func (user User) YesterdayRegisterUser() int {
	now           := time.Now()
	today         := time.Date( now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local )
	todaySec      := today.Unix()
	yesterdaySec  := todaySec - 24 * 60 * 60
	yesterdayTime := time.Unix( yesterdaySec, 0 )
	todayYMD      := today.Format( "2006-01-02" )
	yesterdayYMD  := yesterdayTime.Format( "2006-01-02" )

  return doCount( yesterdayYMD, todayYMD )
}

func (user User) TodayRegisterUser() int {
	now          := time.Now()
	today        := time.Date( now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local )
	todaySec     := today.Unix()
	tomorrowSec  := todaySec + 24 * 60 * 60
	tomorrowTime := time.Unix( tomorrowSec, 0 )
	todayYMD     := today.Format( "2006-01-02" )
	tomorrowYMD  := tomorrowTime.Format( "2006-01-02" )

	return doCount( todayYMD, tomorrowYMD )
}

func doCount( startTime string, endTime string ) int {
  var result int
  var err = DB.Model( &User{} ).Where( "created_at >= ? AND created_at < ?", startTime, endTime ).Count( &result ).Error
  if err != nil {
    return 0
  }
  return result
}
