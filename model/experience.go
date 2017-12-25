package model

type Experience struct {
  ID        uint `gorm:"primary_key" json:"id"`
  Post      uint `json:"post"`
  Level     uint `json:"level"`
  Threshold uint `json:"threshold"`
}

func FetchAllThresholds() []Experience {
  var result []Experience
  DB.Find( &result )

  return result
}

func FetchThresholdBy( level int, post int ) Experience {
  var result Experience
  DB.Where( "level = ? and post = ?", level, post ).Find( &result )

  return result
}

func CalcLevelByExperience( exp uint, post uint ) uint {
  var expRef = FetchAllThresholds()
  result := uint( 1 )
  for _, value := range expRef {
    if post==value.Post && exp >= value.Threshold {
      if value.Level > result {
        result = value.Level
      }
    }
  }

  return result
}