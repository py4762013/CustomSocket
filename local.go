package customSocket

import (
	"log"
	"net"
)

type LsLocal struct {
	Cipher *cipher
	ListenAddr *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

//新建一个本地端
func NewLsLocal(password string, listenAddr, remoteAddr string) (*LsLocal, error) {
	bsPassword, err := parsePassword(password)
	if err != nil {
		return nil, err
	}
	structListenAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	structRemoteAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		return nil, err
	}
	return &LsLocal{
		Cipher: NewCipher(bsPassword),
		ListenAddr: structListenAddr,
		RemoteAddr: structRemoteAddr,
	}, nil
}

//本地端启动监听,接受来自本机浏览器的连接
func (local *LsLocal) Listen(didListen func(listenAddr net.Addr)) error {
	return ListenSecureTcp(local.ListenAddr, local.Cipher, local.handleConn, didListen)
}

func (local *LsLocal) handleConn(userConn *SecureTCPConn) {
	defer userConn.Close()

	proxyServer, err := DialTCPSecure(local.RemoteAddr, local.Cipher)
	if err != nil {
		log.Println(err)
	}

	defer proxyServer.Close()

	//从proxyServer读取数据发送到localUser
	go func() {
		err := proxyServer.DecodeCopy(userConn)
		if err != nil {
			//在copy的过程中可能会存在网络超时等错误
			userConn.Close()
			proxyServer.Close()
		}
	}()
	//从localUser 发送数据到proxyServer
	userConn.EncodeCopy(proxyServer)
}

