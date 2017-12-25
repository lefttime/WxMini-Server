package model

type Benefit struct {
  ID        uint `gorm:"primary_key" json:"id"`
  Post      uint `json:"post"`
  Level     uint `json:"level"`
  Value     uint `json:"value"`
}

func FetchAllBenefit() []Benefit {
  var result []Benefit
  DB.Find( &result )

  return result
}
