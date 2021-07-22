package models

import (
	"net/http"
	"time"
)

//HTTPPacket HTTP 请求包或者响应包
type HTTPPacket struct {
	Stream    string
	FirstLine string
	URL       string
	Body      string
	Header    http.Header
	Time      time.Time
}
