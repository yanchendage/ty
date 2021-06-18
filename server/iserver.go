package server

type IServer interface {
	Start()
	Stop()
	Run()
	AddRouter(msgID uint32, router IRouter)
}