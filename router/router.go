package router

import (
  "gopkg.in/kataras/iris.v6"
  "github.com/lefttime/MyAssistant/config"
  "github.com/lefttime/MyAssistant/controller/order"
  "github.com/lefttime/MyAssistant/controller/product"
  "github.com/lefttime/MyAssistant/controller/promoter"
  "github.com/lefttime/MyAssistant/controller/topup"
  "github.com/lefttime/MyAssistant/controller/user"
)

func Route( app *iris.Framework ) {
  apiPrefix := config.APIConfig.Prefix

  router := app.Party( apiPrefix )
  {
    router.Get( "/wxAppLogin",            user.WxAppLogin            )

    router.Post( "/fetchUserInfo",        user.FetchUserInfo         )
    router.Post( "/setWxAppUser",         user.SetWxAppUserInfo      )
    router.Post( "/searchUserInfo",       user.SearchUserInfo        )

    router.Post( "/payment",              user.Payment               )
    router.Post( "/topupDiamondForUser",  topup.TopupDiamondForUser  )

    router.Post( "/fetchDiamondsInfo",    product.FetchDiamondsInfo  )
    router.Post( "/fetchTopupRecentInfo", topup.FetchTopupRecentInfo )
    router.Post( "/fetchPromoterInfo",    promoter.FetchPromoterInfo )
    router.Post( "/fetchOrdersInfo",      order.FetchOrdersInfo      )
  }
}
