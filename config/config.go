package config

import (
  "fmt"
  "strings"
  "encoding/json"
  "io/ioutil"
  "regexp"
  "github.com/lefttime/MyAssistant/utils"
)

var jsonData map[string]interface{}

func initJSON() {
  bytes, err := ioutil.ReadFile( "./configuration.json" )
  if err != nil {
    fmt.Println( "ReadFile: ", err.Error() )
  }

  configStr := string(bytes[:])
  reg       := regexp.MustCompile( `/\*.*\*/` )

  configStr  = reg.ReplaceAllString( configStr, "" )
  bytes      = []byte( configStr )

  if err := json.Unmarshal( bytes, &jsonData ); err != nil {
    fmt.Println( "invalid config: ", err.Error() )
  }
}

type dBConfig struct {
  Dialect      string
  Database     string
  User         string
  Password     string
  Charset      string
  Host         string
  Port         int
  SQLLog       bool
  URL          string
  MaxIdleConns int
  MaxOpenConns int
}

var DBConfig     dBConfig
var GameDBConfig dBConfig

func initDB() {
  utils.SetStructByJSON( &DBConfig, jsonData["database"].( map[string]interface{} ) )
  portStr := fmt.Sprintf( "%d", DBConfig.Port )
  url     := "{user}:{password}@tcp({host}:{port})/{database}?charset={charset}&parseTime=True&loc=Local"
  url      = strings.Replace( url, "{database}", DBConfig.Database, -1 )
  url      = strings.Replace( url, "{user}",     DBConfig.User,     -1 )
  url      = strings.Replace( url, "{password}", DBConfig.Password, -1 )
  url      = strings.Replace( url, "{host}",     DBConfig.Host,     -1 )
  url      = strings.Replace( url, "{port}",     portStr,           -1 )
  url      = strings.Replace( url, "{charset}",  DBConfig.Charset,  -1 )
  DBConfig.URL = url
}

func initGameDB() {
  utils.SetStructByJSON( &GameDBConfig, jsonData["gameserver"].( map[string]interface{} ) )
  portStr := fmt.Sprintf( "%d", GameDBConfig.Port )
  url     := "{user}:{password}@tcp({host}:{port})/{database}?charset={charset}&parseTime=True&loc=Local"
  url      = strings.Replace( url, "{database}", GameDBConfig.Database, -1 )
  url      = strings.Replace( url, "{user}",     GameDBConfig.User,     -1 )
  url      = strings.Replace( url, "{password}", GameDBConfig.Password, -1 )
  url      = strings.Replace( url, "{host}",     GameDBConfig.Host,     -1 )
  url      = strings.Replace( url, "{port}",     portStr,               -1 )
  url      = strings.Replace( url, "{charset}",  GameDBConfig.Charset,  -1 )
  GameDBConfig.URL = url
}

type serverConfig struct {
  Debug               bool
  ImgPath             string
  UploadImgDir        string
  Port                int
  SessionID           string
  MaxOrder            int
  MinOrder            int
  PageSize            int
  MaxPageSize         int
  MinPageSize         int
  MaxNameLen          int
  MaxRemarkLen        int
  MaxContentLen       int
  MaxProductCateCount int
  MaxProductImgCount  int
}

var ServerConfig serverConfig

func initServer() {
  utils.SetStructByJSON( &ServerConfig, jsonData["go"].( map[string]interface{} ) )
}

type wxAppConfig struct {
  CodeToSessURL string
  AppID         string
  Secret        string
}

var WxAppConfig wxAppConfig

func initWxAppConfig() {
  utils.SetStructByJSON( &WxAppConfig, jsonData["wxApp"].( map[string]interface{} ) )
  url := WxAppConfig.CodeToSessURL
  url  = strings.Replace( url, "{appid}",  WxAppConfig.AppID,  -1 )
  url  = strings.Replace( url, "{secret}", WxAppConfig.Secret, -1 )
  WxAppConfig.CodeToSessURL = url
}

type apiConfig struct {
  Prefix string
  URL    string
}

var APIConfig apiConfig

func initAPI() {
  utils.SetStructByJSON( &APIConfig, jsonData["api"].( map[string]interface{} ) )
}

func init() {
  initJSON()
  initDB()
  initGameDB()
  initServer()
  initWxAppConfig()
  initAPI()
}
