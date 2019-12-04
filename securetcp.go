package customSocket

import (
	"io"
	"log"
	"net"
)

const (
	BufSize = 1024
)

//加密传输的TCP socket
type SecureTCPConn struct {
	io.ReadWriteCloser
	Cipher *cipher
}

//从流中读取加密的数据，解密后把原数据放到bs里
func (secureScoket *SecureTCPConn) DecodeRead(bs []byte) (n int, err error) {
	n, err = secureScoket.Read(bs)
	if err != nil {
		return
	}
	secureScoket.Cipher.decode(bs[:n])
	return
}

//将bs里的数据加密后写入输出流
func (secureSocket *SecureTCPConn) EncodeWrite(bs []byte) (int, error) {
	secureSocket.Cipher.encode(bs)
	return secureSocket.Write(bs)
}

//从src中持续读取原数据加密后写入dst,直到src中没有数据可再读取
func (secureSocket *SecureTCPConn) EncodeCopy(dst io.ReadWriteCloser) error{
	buf := make([]byte, BufSize)
	for {
		readCount, errRead := secureSocket.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, errWrite := (&SecureTCPConn{
				ReadWriteCloser: dst,
				Cipher: secureSocket.Cipher,
			}).EncodeWrite(buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}
			if writeCount != readCount {
				return io.ErrShortWrite
			}
		}
	}
}

//从src中持续读取加密的数据，解密后写入dst,知道src中没有数据可再读取
func (secureSocket *SecureTCPConn) DecodeCopy(dst io.ReadWriteCloser) error{
	buf := make([]byte, BufSize)
	for {
		readCount, errRead := secureSocket.DecodeRead(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, errWrite := dst.Write(buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

// see net.DialTCP
func DialTCPSecure(raddr *net.TCPAddr, cipher *cipher) (*SecureTCPConn, error)  {
	remoteConn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, err
	}
	return &SecureTCPConn{
		ReadWriteCloser: remoteConn,
		Cipher: cipher,
	}, nil
}

// see net.ListenTCp
func ListenSecureTcp(laddr *net.TCPAddr, cipher *cipher, handleConn func(localConn *SecureTCPConn), didListen func(listenAddr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}

	defer listener.Close()

	if listener != nil {
		didListen(listener.Addr())
	}

	for {
		localConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		//当localConn被关闭时，清除所有数据
		localConn.SetLinger(0)
		go handleConn(&SecureTCPConn{
			ReadWriteCloser: localConn,
			Cipher:cipher,
		})
	}
}
