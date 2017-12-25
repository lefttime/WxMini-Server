package model

import (
  "time"
)

type Agency struct {
  ID        uint      `gorm:"primary_key" json:"id"`
  UserID    uint      `json:"userid"`
  Level     uint      `json:"level"`
  LeaderID  uint      `json:"leaderid"`
  CreatedAt time.Time `json:"createdAt"`
}

func UpdateAgencyInfo( agency Agency ) {
  var count int
  DB.Model( &Agency{} ).Where( map[string]interface{}{ "userid": agency.UserID, "level": agency.Level, "leader_id": agency.LeaderID } ).Count( &count )
  if count==0 {
    agency.CreatedAt = time.Now()
    DB.Create( &agency )
  }
}

func FetchAgencysByLeaderID( leaderID int ) []Agency {
  var result []Agency
  DB.Where( "leader_id = ?", leaderID ).Find( &result )

  return result
}