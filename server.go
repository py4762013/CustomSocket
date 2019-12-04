package customSocket

import (
	"encoding/binary"
	"net"
)

type LsServer struct {
	Cipher *cipher
	ListenAddr *net.TCPAddr
}

//新建服务器端
func NewLsServer(password string, listenAddr string) (*LsServer, error) {
	bsPassword, err := parsePassword(password)
	if err != nil {
		return nil, err
	}
	structListenAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	return &LsServer{
		Cipher: NewCipher(bsPassword),
		ListenAddr: structListenAddr,
	}, nil
}

func (LsServer *LsServer) Listen(didListen func(listenAddr net.Addr)) error {
	return ListenSecureTcp(LsServer.ListenAddr, LsServer.Cipher, LsServer.handleConn, didListen)
}

//解SOCKS5协议
func (lsServer *LsServer) handleConn(localConn *SecureTCPConn) {
	defer localConn.Close()
	buf := make([]byte, 256)

	/**
	The localConn connects to dtServer, and sends a ver identifier/method selection message:
	the Ver field set to X'05' for this ver of the protocol.
	              +----+----------+----------+
		          |VER | NMETHODS | METHODS  |
		          +----+----------+----------+
		          | 1  |    1     | 1 to 255 |
		          +----+----------+----------+
	 */

	_, err := localConn.DecodeRead(buf)

	if err != nil || buf[0] != 0x05 {
		return
	}

	/**
			  +----+--------+
	          |VER | METHOD |
	          +----+--------+
	          | 1  |   1    |
	          +----+--------+
	 */

	// 不需要验证，直接校验通过
	localConn.EncodeWrite([]byte{0x05, 0x00});

	/**
	  +----+-----+-------+------+----------+----------+
	  |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
	*/
	// 获取真正的远程服务地址
	n, err := localConn.DecodeRead(buf)
	// n 的最短长度为7
	if err != nil || n < 7 {
		return
	}

	// cmd代表客户端的请求类型,值的长度为1个字节，类型有三种
	if buf[1] != 0x01 {
		return
	}

	var dIP []byte
	//aType代表请求的远程服务器地址类型,值长度1字段，有三种类型
	switch buf[3] {
	case 0x01:
		//IP V4 address； X'01'
		dIP = buf[4 : 4+net.IPv4len]
	case 0x03:
		// Domain Name X'03'
		ipAddr, err := net.ResolveIPAddr("ip", string(buf[5:n-2]))
		if err != nil {
			return
		}
		dIP = ipAddr.IP
	case 0x04:
		// IP v6 address: X'04'\
		dIP = buf[4: 4+net.IPv6len]
	default:
		return
	}
	dPort := buf[n-2:]
	dstAddr := &net.TCPAddr{
		IP: dIP,
		Port: int(binary.BigEndian.Uint16(dPort)),
	}

	//连接真正的服务
	dstServer, err := net.DialTCP("tcp",nil, dstAddr)
	if err != nil {
		return
	} else {
		defer dstServer.Close()
		// Conn被关闭的时候 清楚所有数据
		dstServer.SetLinger(0)

		//响应客户端连接成功
		localConn.EncodeWrite([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00});
	}

	//进行转发
	//从localUser读取数据转发到dstServer
	go func() {
		err := localConn.DecodeCopy(dstServer)
		if err != nil {
			//在copy中出现错误
			localConn.Close()
			dstServer.Close()
		}
	}()
	//从dstSever 读取数据发送到localUser
	(&SecureTCPConn{
		Cipher: localConn.Cipher,
		ReadWriteCloser: dstServer,
	}).EncodeCopy(localConn)
}