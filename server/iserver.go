package server

type IServer interface {
	Start()
	Stop()
	Run()
	AddRouter(msgID uint32, router IRouter)

	SetOnServerStartCallback(func() error)
	SetOnServerStopCallback(func() error)
	CallOnServerStartCallback() error
	CallOnServerStopCallback() error

	SetOnConnStartCallback(func (IConnection))
	SetOnConnStopCallback(func (IConnection))
	CallOnConnStartCallback(conn IConnection)
	CallOnConnStopCallback(conn IConnection)

	SetProperty(key string, value interface{})
	GetProperty(key string)(interface{}, error)
	RemoveProperty(key string)
}