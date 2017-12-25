package model

type Discount struct {
  ID        uint `gorm:"primary_key" json:"id"`
  ProductID uint `json:"productId"`
  Post      uint `json:"post"`
  Level     uint `json:"level"`
  Value     uint `json:"value"`
}

func FetchAllDiscount() []Discount {
  var result []Discount
  DB.Find( &result )

  return result
}