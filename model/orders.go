package model

import (
  "time"
)

type Order struct {
  ID           uint      `gorm:"primary_key" json:"id"`
  UserID       uint      `json:"userId"`
  ProductID    uint      `json:"productId"`
  ProductCount uint      `json:"productCount"`
  Rebate       uint      `json:"rebate"`
  Discount     uint      `json:"discount"`
  OriginPrice  float32   `json:"originPrice"`
  TotalPrice   float32   `json:"totalPrice"`
  Status       uint      `json:"status"`
  CreatedAt    time.Time `json:"createdAt"`
  UpdatedAt    time.Time `json:"updatedAt"`
  PayAt        time.Time `json:"payAt"`
}

func FetchOrdersByUserId( userId int ) []Order {
  var result []Order
  DB.Where( "user_id = ?", userId ).Find( &result )

  return result
}
