package model

type Product struct {
  ID    uint `gorm:"primary_key" json:"id"`
  Type  uint `json:"type"`
  Count uint `json:"count"`
  Value uint `json:"value"`
}

func FetchProductsByType( vType int ) []Product {
  var result []Product
  DB.Where( "type = ?", vType ).Find( &result )

  return result
}
