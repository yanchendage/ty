package server

import "net"

type IConnection interface {
	Start()
	Stop()

	GetServer() IServer
	GetTCPConn() *net.TCPConn
	GetConnID() uint32
	GetRemoteAddr() net.Addr

	SendMsg(msgID uint32, data []byte) error
	SendBuffMsg(msgID uint32, data []byte) error

	SetProperty(key string, value interface{})
	GetProperty(key string)(interface{}, error)
	RemoveProperty(key string)
}

type Handle func(*net.TCPConn, []byte, int) error

