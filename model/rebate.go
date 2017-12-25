package model

type Rebate struct {
  ID        uint `gorm:"primary_key" json:"id"`
  Post      uint `json:"post"`
  Level     uint `json:"level"`
  Value     uint `json:"value"`
}

func FetchAllRebates() []Rebate {
  var result []Rebate
  DB.Find( &result )

  return result
}