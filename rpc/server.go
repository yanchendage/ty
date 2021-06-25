package rpc

import (
	"errors"
	"github.com/yanchendage/ty/server"
	"log"
	"strings"
	"sync"
)

type GobCoderRouter struct {
	server.BRouter
}

func (gc *GobCoderRouter) Handle(request server.IRequest) {

	log.Println("【RPC Server】 request msgId=", request.GetMsgID(), ", data=", string(request.GetMsgData()))

	//coder
	f := NewCoderFuncMap[GobProtocol]
	if f == nil {
		log.Println("【RPC Server】invalid coder type ", request.GetMsgID())
		return
	}

	coder := f(request.GetConnection().GetTCPConn())

	//read msg data
	msg, err := coder.Decode(request.GetMsgData())
	if err !=nil {
		log.Println("【RPC Server】coder decode err ", err)
		return
	}

	serviceManager,_ := request.GetConnection().GetServer().GetProperty("ServiceManager")
	sm := serviceManager.(*ServiceManager)

	svc , mt, err := sm.findService(msg.H.ServiceMethod)
	if err != nil {
		log.Println("【RPC Server】find service method err ", err)
		return
	}

	//msg.methodType = mt
	//msg.svc = svc

	//arg := msg.methodType.newArgv()
	//reply := msg.methodType.newReplyv()

	log.Println(svc,mt)

	//arg.Set(msg.B.Args)
	//
	//svc.call(mt,arg,reply)
	//
	//log.Println(*reply.Interface().(*int))

	//msg.B.Arg = reflect.New(reflect.TypeOf(""))
	//
	//if err := coder.ReadBody(req.body.Arg); err != nil{
	//	log.Println("rpc server: read body error:", err)
	//}
	//
	//log.Println(req)
	//log.Println(req.body.Arg.Elem())

	//req.replyv = reflect.ValueOf(fmt.Sprintf("geerpc resp %d", req.h.Seq))
	//
	//
	//request.GetConnection().SendMsg(0, []byte("pong"))

	//handleRequest


	//sendResponse

	//err := request.GetConnection().SendMsg(0, []byte("pong"))
}

type ServiceManager struct {
	ServiceMap sync.Map
}

func (sm *ServiceManager) RegisterService(service interface{}) error {
	s := newService(service)
	if _ , loaded := sm.ServiceMap.LoadOrStore(s.name, s); loaded {
		return errors.New("rpc: service already defined: " + s.name)
	}
	log.Println("【RPC Server】register service ")
	return  nil
}

func (sm *ServiceManager) findService(serviceMethod string) (svc *service, methodType *methodType, err error) {
	index := strings.LastIndex(serviceMethod, ".")

	if index < 0 {
		err = errors.New("invalid serviceMethod ")
	}
	serviceName, methodName := serviceMethod[:index], serviceMethod[index+1:]
	svcInterface, ok := sm.ServiceMap.Load(serviceName)

	if !ok {
		err = errors.New("rpc server: can't find service " + serviceName)
		return
	}

	svc = svcInterface.(*service)
	methodType = svc.method[methodName]

	if methodType == nil{
		err = errors.New("rpc server: can't find method " + methodName)
	}
	return
}



func NewServiceManager() *ServiceManager {
	return &ServiceManager{}
}

func NewServerManager(serverName string, host string, port int,properties map[string]interface{}){

	r := server.NewServer(serverName, host, port)

	for key, value := range properties {
		r.SetProperty(key, value)
	}

	r.AddRouter(0,&GobCoderRouter{})

	r.Run()
}

func NewServer(serverName string, host string, port int)  {
	r := server.NewServer(serverName, host, port)
	r.AddRouter(0,&GobCoderRouter{})
	r.Run()
}



