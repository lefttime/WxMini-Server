package model

type errorCode struct {
  Success      int
  Error        int
  SessionError int
  NotFound     int
}

var ErrorCode = errorCode {
  Success      : 0,
  Error        : 1,
  SessionError : 400,
  NotFound     : 404,
}