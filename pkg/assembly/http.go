package assembly

import (
	"bufio"
	"fmt"
	"log"
	"net/http"

	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type httpStreamFactory struct{}

type httpStream struct {
	netLayer       gopacket.Flow
	transportLayer gopacket.Flow
	reader         *tcpreader.ReaderStream
}

func (h *httpStream) run() {
	//1.将reader 转化为一个带buffer 的reader
	readerWithBuffer := bufio.NewReader(h.reader)
	for {
		request, err := http.ReadRequest(readerWithBuffer)
		if err != nil {
			log.Println("http.ReadRequest failed,err:", err)
			return
		}
		defer request.Body.Close()
		fmt.Println(h.netLayer.String())
	}
}
