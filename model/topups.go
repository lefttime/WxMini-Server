package model

import (
  "time"
)

type Topup struct {
  ID           uint      `gorm:"primary_key" json:"id"`
  UserID       uint      `json:"userId"`
  PresenteeID  uint      `json:"presenteeId"`
  Count        uint      `json:"count"`
  CreatedAt    time.Time `json:"createdAt"`
  Status       int       `json:"status"`                // -1: 失败  0: 等待  1: 成功
}

type TopupPresenteeInfo struct {
  OpenID     string `json:"openId"`
  Avatar     string `json:"avatar"`
  Nickname   string `json:"nickname"`
  Post       uint   `json:"post"`
  GameID     uint   `json:"gameId"`
  Phone      string `json:"phone"`
  Count      uint   `json:"count"`
  CreatedAt  string `json:"datetime"`
  Status     int    `json:"status"`
}

func FetchTopupsByUserId( userId int, limit int ) []Topup {
  var result []Topup
  if limit <= 0 {
    DB.Where( "user_id = ?", userId ).Find( &result )
  } else {
    DB.Limit( limit ).Where( "user_id = ?", userId ).Find( &result )
  }

  return result
}

func FetchTopupsByPresenteeID( presenteeId int ) []Topup {
  var result []Topup
  DB.Where( "presentee_id = ?", presenteeId ).Find( &result )

  return result
}

func FetchUserTodayTopups( userId int ) []Topup {
  now         := time.Now()
  today       := time.Date( now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local )
  todaySec    := today.Unix()
  tomorrowSec := todaySec + 24 * 60 * 60
  tomorrow    := time.Unix( tomorrowSec, 0 )
  return FetchUserTopupsBetween( today, tomorrow, userId )
}

func FetchUserLastWeekTopups( userId int ) []Topup {
  lastMonday, lastSunday := lastWeekRange()
  return FetchUserTopupsBetween( lastMonday, lastSunday, userId )
}

func FetchUserTopupsBetween( startTime time.Time, stopTime time.Time, userId int ) []Topup {
  var result []Topup

  startYMD := FormatDatetime( startTime, true )
  stopYMD  := FormatDatetime( stopTime,  true  )
  DB.Where( "user_id = ? and created_at >= ? and created_at <= ?", userId, startYMD, stopYMD ).Find( &result )

  return result
}

func FormatDatetime( resDatetime time.Time, short bool ) string {
  if short {
    return resDatetime.Format( "2006-01-02" )
  }
  return resDatetime.Format( "2006-01-02 15:04:05" )
}

func lastWeekRange() (lastMonday time.Time, lastSunday time.Time) {
  now      := time.Now()
  today    := time.Date( now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local )
  todaySec := today.Unix()

  duration := int64( now.Weekday() )
  if duration==0 {
    duration = 13
  } else {
    duration = duration + 6
  }
  lastMonSec := todaySec - duration * 24 * 60 * 60
  lastSunSec := lastMonSec + 7 * 24 * 60 * 60

  return time.Unix( lastMonSec, 0 ), time.Unix( lastSunSec, 0 )
}
