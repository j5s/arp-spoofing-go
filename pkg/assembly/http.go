package assembly

import (
	"ARPSpoofing/dao/memory"
	"ARPSpoofing/models"
	"time"

	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

func getProtocal(dstport string) string {
	if dstport == "443" {
		return "https"
	}
	return "http"
}

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
			//0.读取网络层中有用的数据
			stream := fmt.Sprintf("%s:%s->%s:%s\n", netLayer.Src().String(), transportLayer.Src().String(), netLayer.Dst().String(), transportLayer.Dst().String())
			url := fmt.Sprintf("%s://%s%s\n\n", getProtocal(transportLayer.Dst().String()), request.Host, request.RequestURI)
			//http 头首行
			firstLine := fmt.Sprintf("%s %s %s\n", request.Method, request.RequestURI, request.Proto)
			//http 体
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				continue
			}
			memory.DataCh <- &models.HTTPPacket{
				Stream:    stream,
				URL:       url,
				FirstLine: firstLine,
				Header:    request.Header,
				Body:      string(body),
				Time:      time.Now(),
			}
		}
		//2.尝试从流中读取http响应
		response, err := http.ReadResponse(streamReaderWithBuffer, nil)
		if err == nil {
			defer response.Body.Close()
			//0.读取网络层中有用的数据
			stream := fmt.Sprintf("%s:%s->%s:%s\n", netLayer.Src().String(), transportLayer.Src().String(), netLayer.Dst().String(), transportLayer.Dst().String())
			//1.http 体
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				continue
			}
			memory.DataCh <- &models.HTTPPacket{
				Stream:    stream,
				FirstLine: fmt.Sprintf("%s %d %s", response.Proto, response.StatusCode, response.Status),
				Header:    response.Header,
				Body:      string(body),
				Time:      time.Now(),
			}
		}
		//3.尝试从流中读取https 报文
		if transportLayer.Dst().String() == "443" {
			temp, err := ioutil.ReadAll(streamReaderWithBuffer)
			if err == nil && len(temp) > 0 {
				memory.DataCh <- &models.HTTPPacket{
					Stream: fmt.Sprintf("%s:%s->%s:%s\n", netLayer.Src().String(), transportLayer.Src().String(), netLayer.Dst().String(), transportLayer.Dst().String()),
					Body:   string(temp),
					Time:   time.Now(),
				}
			}
		}
	}
}
