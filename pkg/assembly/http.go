package assembly

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

//httpStreamFactory 流工厂
type httpStreamFactory struct{}

//New 每来一个tcp会话，就调用流工厂的New方法创建一个流加入到流池中，然后会调用stream的重组方法，将tcp数据包重组到流，要读取重组后的结果得创建一个协程不断读取流
func (f *httpStreamFactory) New(netLayer, transportLayer gopacket.Flow) tcpassembly.Stream {
	//实现了 Stream 接口和 Reader接口
	stream := tcpreader.NewReaderStream()
	go ReadFromStream(&stream, netLayer, transportLayer)
	return &stream
}

//ReadFromStream 从流中读取数据
func ReadFromStream(streamReader io.Reader, netLayer, transportLayer gopacket.Flow) {
	//1.将reader 转化为一个带buffer 的reader,因为http.ReadRequest 需要
	streamReaderWithBuffer := bufio.NewReader(streamReader)
	for {
		//1.尝试从流中读取http请求（流中不一定是http请求）
		request, err := http.ReadRequest(streamReaderWithBuffer)
		if err == nil {
			defer request.Body.Close()
			//输出http 头
			for key, values := range request.Header {
				fmt.Printf("%v:%v\n", key, strings.Join(values, ""))
			}
			fmt.Println()
			//输出http 体
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				continue
			}
			fmt.Println(string(body))
			fmt.Printf("\n\n")
		}
		//2.尝试从流中读取http响应
		response, err := http.ReadResponse(streamReaderWithBuffer, nil)
		if err == nil {
			for key, values := range response.Header {
				fmt.Printf("%v:%v\n", key, strings.Join(values, ""))
			}
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				continue
			}
			fmt.Println(string(body))
			fmt.Printf("\n\n")
		}
	}
}
