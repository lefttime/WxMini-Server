package rebate

import (
//   "encoding/json"
//   "gopkg.in/kataras/iris.v6"
  "github.com/lefttime/MyAssistant/model"
//   "github.com/lefttime/MyAssistant/controller/common"
)

func FetchRebatesInfo() []model.Rebate {
  return model.FetchAllRebates()
}

func CalcRebate( post uint, level uint, productCount uint ) uint {
  var rebates = FetchRebatesInfo()
  for _, value := range rebates {
    if value.Post==post && value.Level==level {
      return value.Value * productCount / 100
    }
  }

  return 0
}