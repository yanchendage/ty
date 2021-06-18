package server

type IMessageManager interface {
	AddRouter(msgID uint32, router IRouter)
	GetRouter(msgID uint32) IRouter
	StartWorkerPool()
	AddWorkerToWorkQueue(request IRequest)
}