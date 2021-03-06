package main

import (
  "fmt"
  "os"
  "time"
  "strconv"

  _ "github.com/jinzhu/gorm/dialects/mysql"

  "github.com/jinzhu/gorm"
  "gopkg.in/kataras/iris.v6"
  "gopkg.in/kataras/iris.v6/adaptors/httprouter"
  "gopkg.in/kataras/iris.v6/adaptors/sessions"

  "github.com/lefttime/MyAssistant/config"
  "github.com/lefttime/MyAssistant/model"
  "github.com/lefttime/MyAssistant/router"
)

func init() {
  db, err         := gorm.Open( config.DBConfig.Dialect,     config.DBConfig.URL     )
  gameDb, gameErr := gorm.Open( config.GameDBConfig.Dialect, config.GameDBConfig.URL )
  if err != nil || gameErr != nil {
    fmt.Println( err.Error() )
    fmt.Println( gameErr.Error() )
    os.Exit( -1 )
  }

  if config.DBConfig.SQLLog {
    db.LogMode( true )
  }

  if config.GameDBConfig.SQLLog {
    gameDb.LogMode( true )
  }

  db.DB().SetMaxIdleConns( config.DBConfig.MaxIdleConns )
  db.DB().SetMaxOpenConns( config.DBConfig.MaxOpenConns )

  gameDb.DB().SetMaxIdleConns( config.GameDBConfig.MaxIdleConns )
  gameDb.DB().SetMaxOpenConns( config.GameDBConfig.MaxOpenConns )

  model.DB     = db
  model.GameDB = gameDb
}

func main() {
  app := iris.New( iris.Configuration {
    Gzip    : true,
    Charset : "UTF-8",
  })

  if config.ServerConfig.Debug {
    app.Adapt( iris.DevLogger() )
  }

  app.Adapt( sessions.New( sessions.Config {
    Cookie: config.ServerConfig.SessionID,
    Expires: time.Minute * 20,
  }))

  app.Adapt( httprouter.New() )
  router.Route( app )

  app.OnError( iris.StatusNotFound, func( ctx *iris.Context ) {
    ctx.JSON( iris.StatusOK, iris.Map {
      "errNo" : model.ErrorCode.NotFound,
      "msg"   : "Not Found",
      "data"  : iris.Map{},
    } )
  } )

  app.OnError( 500, func( ctx *iris.Context ) {
    ctx.JSON( iris.StatusInternalServerError, iris.Map {
      "errNo" : model.ErrorCode.Error,
      "msg"   : "error",
      "data"  : iris.Map{},
    } )
  } )

  app.Listen( ":" + strconv.Itoa( config.ServerConfig.Port ) )
}