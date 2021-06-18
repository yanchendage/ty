package server

import (
	"fmt"
	"log"
	"testing"
)


type PongRouter struct {
	BRouter
}

func (this *PongRouter) Handle(request IRequest) {
	log.Println("【Request】 msgId=", request.GetMsgID(), ", data=", string(request.GetMsgData()))

	err := request.GetConnection().SendMsg(0, []byte("pong"))
	if err != nil {
		fmt.Println(err)
	}
}

//HelloZinxRouter Handle
type HelloZinxRouter struct {
	BRouter
}

func (this *HelloZinxRouter) Handle(request IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetMsgData()))

	err := request.GetConnection().SendMsg(1, []byte("Hello Zinx Router V0.6"))
	if err != nil {
		fmt.Println(err)
	}
}

func TestNewServer(t *testing.T) {
	s := NewServer()
	//配置路由
	s.AddRouter(0, &PongRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	s.Run()
}
