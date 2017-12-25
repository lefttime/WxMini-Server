package model

import (
  "time"
)

type Promoter struct {
  ID        uint      `gorm:"primary_key" json:"id"`
  UserID    uint      `json:"userid"`
  Level     uint      `json:"level"`
  LeaderID  uint      `json:"leaderid"`
  CreatedAt time.Time `json:"createdAt"`
}

type PromoterInfo struct {
  OpenID   string `json:"openId"`
  Avatar   string `json:"avatar"`
  Nickname string `json:"nickname"`
  Diamond   int   `json:"diamond"`
  Phone    string `json:"phone"`
}

func UpdatePromoterInfo( promoter Promoter ) {
  var count int
  DB.Model( &Promoter{} ).Where( map[string]interface{}{ "userid": promoter.UserID, "level": promoter.Level, "leader_id": promoter.LeaderID } ).Count( &count )
  if count==0 {
    promoter.CreatedAt = time.Now()
    DB.Create( &promoter )
  }
}

func FetchPromotersByLeaderID( leaderID int, limit int ) []Promoter {
  var result []Promoter
  if limit <= 0 {
    DB.Where( "leader_id = ?", leaderID ).Find( &result )
  } else {
    DB.Limit( limit ).Where( "leader_id = ?", leaderID ).Find( &result )
  }

  return result
}