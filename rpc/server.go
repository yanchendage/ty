package rpc

import (
	"errors"
	"fmt"
	"github.com/yanchendage/ty/server"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"
)


//
//
//type GobCoderRouter struct {
//	server.BRouter
//}
//
//func (gc *GobCoderRouter) Handle(request server.IRequest) {
//
//	log.Println("【RPC Server】 request msgId=", request.GetMsgID(), ", data=", string(request.GetMsgData()))
//
//	//new coder
//	f := NewCoderFuncMap[GobProtocol]
//	if f == nil {
//		log.Println("【RPC Server】invalid coder type ", request.GetMsgID())
//		return
//	}
//
//	coder := f(request.GetConnection().GetTCPConn())
//
//	header, err := coder.DecodeHeader(request.GetMsgData())
//
//	if err !=nil {
//		log.Println("【RPC Server】coder decode header err ", err)
//		return
//	}
//
//	serviceManager,_ := request.GetConnection().GetServer().GetProperty("ServiceManager")
//	sm := serviceManager.(*ServiceManager)
//
//	svc, mt, err := sm.findService(header.ServiceMethod)
//
//	if err != nil {
//		log.Println("【RPC Server】find service method err ", err)
//		return
//	}
//
//	arg := mt.newArgv()
//	reply := mt.newReplyv()
//	//
//	argvi := arg.Interface()
//	if arg.Type().Kind() != reflect.Ptr {
//		argvi = arg.Addr().Interface()
//	}
//
//	arg, err = coder.DecodeArgs(request.GetMsgData(),argvi)
//	if err != nil {
//		log.Println("【RPC Server】Decode Args err ", err)
//		header.Err = "Decode Args err"
//		return
//	}
//
//	err = svc.call(mt, arg, reply)
//	if err != nil {
//		log.Println("【RPC Server】service call err ", err)
//		header.Err = "service call err"
//		return
//	}
//
//	msg := Msg{
//		H: header,
//		Reply:reply,
//	}
//
//	respBuf, err := coder.Encode(msg)
//
//	if err != nil {
//		log.Println("【RPC Server】response coder encode err ", err)
//		return
//	}
//
//	err = request.GetConnection().SendMsg(0,respBuf)
//
//	if err != nil {
//		log.Println("【RPC Server】response send msg err ", err)
//		return
//	}
//}

type ProtoCoderRouter struct {
	server.BRouter
}

func (gc *ProtoCoderRouter) Handle(request server.IRequest) {
	//new coder
	f := NewCoderFuncMap[ProtobufProtocol]
	if f == nil {
		log.Println("【RPC Server】invalid coder type ", request.GetMsgID())
		return
	}

	coder := f(request.GetConnection().GetTCPConn())

	header, err := coder.DecodeHeader(request.GetMsgData())

	if err !=nil {
		log.Println("【RPC Server】coder decode header err ", err)
		return
	}
	
	serviceManager,_ := request.GetConnection().GetServer().GetProperty("ServiceManager")
	sm := serviceManager.(*ServiceManager)

	svc, mt, err := sm.findService(header.ServiceMethod)
	if err != nil {
		log.Println("【RPC Server】find service method err ", err)
		return
	}

	arg := mt.newArgv()
	reply := mt.newReplyv()
	//
	argvi := arg.Interface()
	if arg.Type().Kind() != reflect.Ptr {
		argvi = arg.Addr().Interface()
	}

	arg, err = coder.DecodeBody(request.GetMsgData(),argvi)
	if err != nil {
		log.Println("【RPC Server】Decode Args err ", err)
		header.Err = "Decode Args err"
		return
	}

	err = svc.call(mt, arg, reply)
	if err != nil {
		log.Println("【RPC Server】service call err ", err)
		header.Err = "service call err"
		return
	}

	respBuf, err := coder.EncodeResponse(&header, reply.Interface())
	if err != nil {
		log.Println("【RPC Server】response coder encode err ", err)
		return
	}

	err = request.GetConnection().SendMsg(1,respBuf)

	if err != nil {
		log.Println("【RPC Server】response send msg err ", err)
		return
	}

	//
	//msg := Msg{
	//	H: header,
	//	Reply:reply,
	//}

	//fmt.Println(msg)

	//respBuf, err := coder.Encode(msg)
	//
	//if err != nil {
	//	log.Println("【RPC Server】response coder encode err ", err)
	//	return
	//}
	//
	//err = request.GetConnection().SendMsg(0,respBuf)
	//
	//if err != nil {
	//	log.Println("【RPC Server】response send msg err ", err)
	//	return
	//}
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
		fmt.Println("invalid serviceMethod")
		err = errors.New("invalid serviceMethod")
	}
	serviceName, methodName := serviceMethod[:index], serviceMethod[index+1:]
	svcInterface, ok := sm.ServiceMap.Load(serviceName)

	if !ok {
		fmt.Println("rpc server: can't find service ")
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


type ServerManager struct {
	serviceManager *ServiceManager
	Server server.IServer
	RegistryAddr string
}

func (serverManager *ServerManager) AddProperty(key string, value interface{})  {
	serverManager.Server.SetProperty(key, value)
}

func (serverManager *ServerManager) RegisterService(service interface{}) error {
	s := newService(service)
	if _ , loaded := serverManager.serviceManager.ServiceMap.LoadOrStore(s.name, s); loaded {
		return errors.New("rpc: service already defined: " + s.name)
	}
	log.Println("【RPC Server】register service ")
	return  nil
}

func (serverManager *ServerManager) Run() {

	serverManager.Server.SetProperty("ServiceManager", serverManager.serviceManager)

	//serverManager.Server.AddRouter(0,&GobCoderRouter{})
	serverManager.Server.AddRouter(1,&ProtoCoderRouter{})

	serverManager.Server.Run()
}


func InitServerManagerAndRegister(serverName string, host string, port int, registryAddr string) *ServerManager{
	r := server.NewServer(serverName, host, port)
	r.SetProperty("registryAddr",registryAddr)

	r.SetOnServerStartCallback(func() error {
		registryAddr,err := r.GetProperty("registryAddr")
		addr := fmt.Sprintf("%s?addr=%s:%d",registryAddr.(string), r.IP, r.Port)

		if err !=nil {
			return err
		}

		_, err = http.Post(addr,
			"application/x-www-form-urlencoded", nil)
		if err != nil {
			log.Println("【server】 register to rpc registry err",err)
			return err
		}

		return err
	})

	return &ServerManager{
		serviceManager: NewServiceManager(),
		Server:         r,
		RegistryAddr: registryAddr,
	}
}


func InitServerManager(serverName string, host string, port int) *ServerManager{
	r := server.NewServer(serverName, host, port)
	return &ServerManager{
		serviceManager: NewServiceManager(),
		Server:         r,
	}
}

func NewServerManager(serverName string, host string, port int,properties map[string]interface{}){

	r := server.NewServer(serverName, host, port)

	for key, value := range properties {
		r.SetProperty(key, value)
	}

	//r.AddRouter(0,&GobCoderRouter{})


	r.Run()

	//service register
	fmt.Println("todo service register")

}

func NewServer(serverName string, host string, port int)  {
	r := server.NewServer(serverName, host, port)
	//r.AddRouter(0,&GobCoderRouter{})
	r.Run()
}



