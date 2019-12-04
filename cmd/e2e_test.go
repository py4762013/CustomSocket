package cmd

import (
	"github.com/py4762013/CustomSocket"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"math/rand"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

const (
	MaxPackSize					= 1024*1024*5 // 5Mb
	EchoServerAddr				= "127.0.0.1:3453"
	CustomSocksProxyLocalAddr	= "127.0.0.1:8448"
	CustomSocksProxyServerAddr	= "127.0.0.1:8449"
)

var (
	customsocksDialer proxy.Dialer
)

func init()  {
	log.SetFlags(log.Lshortfile)
	go runEchoServer()
	go runCustomSocketsProxyServer()
	//初始化代理socksDialer
	var err error
	time.Sleep(time.Second)
	customsocksDialer, err = proxy.SOCKS5("tcp", CustomSocksProxyLocalAddr, nil, proxy.Direct)
	if err != nil {
		log.Fatalln(err)
	}
}

// 启动echo server
func runEchoServer()  {
	listener, err := net.Listen("tcp", EchoServerAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		log.Println("echoServer Connect Accept")
		go func() {
			defer func() {
				conn.Close()
				log.Println("echoServer Connect Close")
			}()
			io.Copy(conn, conn)
		}()
	}
}

func runCustomSocketsProxyServer()  {
	password := customSocket.RandPassword()
	serverS, err := customSocket.NewLsLocal(password, CustomSocksProxyLocalAddr, CustomSocksProxyServerAddr)
	if err != nil {
		log.Fatalln(err)
	}
	localS, err := customSocket.NewLsServer(password, CustomSocksProxyServerAddr)
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		log.Fatalln(serverS.Listen(func(listenAddr net.Addr) {
			log.Println(listenAddr)
		}))
	}()
	log.Fatalln(localS.Listen(func(listenAddr net.Addr) {
		log.Println(listenAddr)
	}))
}

// 发送一次连接测试经过代理扣数据传输的正确性
// packSize 代表这个连接发送数据的大小
func testConnect(packSize int)  {
	//随机生产 MaxPackSize byte的[]byte
	data := make([]byte, packSize)
	_, err := rand.Read(data)

	// 连接
	conn, err := customsocksDialer.Dial("tcp", EchoServerAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	//写
	go func() {
		conn.Write(data)
	}()

	//读
	buf := make([]byte, len(data))
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		log.Fatalln(err)
	}
	if !reflect.DeepEqual(data, buf) {
		log.Fatalln("通过 customSocks 代理传输得到的数据前后不一致")
	}
}

func testCustomSocks(t *testing.T)  {
	testConnect(rand.Intn(MaxPackSize))
}

//获取并发送 data 到 echo server 并且收到全部返回所花费的时间
func benchmarkCustomSocks(concurrenceCount int)  {
	wg := sync.WaitGroup{}
	wg.Add(concurrenceCount)
	for i := 0; i < concurrenceCount; i++ {
		go func() {
			testConnect(rand.Intn(MaxPackSize))
			wg.Done()
		}()
	}
	wg.Wait()
}

//获取 发送 data 到 echo server 并且收到全部返回 所花费的时间
func BenchmarkCustomSocks(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		benchmarkCustomSocks(10)
		b.StopTimer()
	}
}