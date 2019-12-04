package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/phayes/freeport"
	"github.com/py4762013/CustomSocket"
	"github.com/py4762013/CustomSocket/cmd"
)

var version ="master"

func main()  {
	log.SetFlags(log.Lshortfile)

	//优先从环境变量中获取监听端口
	port, err := strconv.Atoi(os.Getenv("CUSTOMSOCKS_SERVER_PORT"))
	//服务监听端口随机生成
	if err != nil {
		port, err = freeport.GetFreePort()
	}
	if err != nil {
		//随机端口失败就采用7448
		port = 7448
	}
	//默认配置
	config := &cmd.Config{
		ListenAddr: fmt.Sprintf(":%d", port),
		Password:   customSocket.RandPassword(),
	}
	config.ReadConfig()
	config.SaveConfig()

	//启动server端并监听
	lsServer, err := customSocket.NewLsServer(config.Password, config.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatalln(lsServer.Listen(func(listenAddr net.Addr) {
		log.Println(fmt.Sprintf(`
customsocks-server:%s 启动成功，配置如下：
服务器监听地址：
%s
密码:
%s`, version, listenAddr, config.Password))
	}))
}
