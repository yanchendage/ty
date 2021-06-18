package server

import "net"

type IConnection interface {
	Start()
	Stop()

	GetTCPConn() *net.TCPConn
	GetConnID() uint32
	GetRemoteAddr() net.Addr

	SendMsg(msgID uint32, data []byte) error
	SendBuffMsg(msgID uint32, data []byte) error
}

type Handle func(*net.TCPConn, []byte, int) error

