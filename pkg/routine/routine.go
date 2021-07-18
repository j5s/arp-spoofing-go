package routine

import (
	"bytes"
	"log"
	"runtime"
	"strconv"
)

//GetGID 获取协程id
func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		log.Println(err)
	}
	return n
}
