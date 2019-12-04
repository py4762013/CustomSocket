package main

import (
	"fmt"
	"log"
	"net"

	"github.com/py4762013/CustomSocket"
	"github.com/py4762013/CustomSocket/cmd"
)

const (
	DefaultListenAddr = ":7448"
)

var version = "master"

func main()  {
	log.SetFlags(log.Lshortfile)

	//默认配置
	config := &cmd.Config{
		ListenAddr:DefaultListenAddr,
	}
	config.ReadConfig()
	config.SaveConfig()

	//启动local端并监听
	lsLocal, err := customSocket.NewLsLocal(config.Password, config.ListenAddr, config.RemoteAddr)
	if err != nil {
		log.Println(err)
	}
	log.Fatalln(lsLocal.Listen(func(listenAddr net.Addr) {
		log.Println(fmt.Sprintf(`
customsocks-local: %s 启动成功,配置如下
本地监听地址:
%s
远程服务地址:
%s
密码:
%s`, version, listenAddr, config.RemoteAddr, config.Password))
	}))
}
