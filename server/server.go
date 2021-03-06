package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	Name string
	IP string
	Port int
	IPVersion string
	MessageManager IMessageManager
	ConnectionManager IConnectionManager

	OnServerStart func() error
	OnServerStop func() error

	OnConnStart	func(conn IConnection)
	OnConnStop func(conn IConnection)

	property map[string]interface{}
	propertyLock sync.RWMutex

	stop chan bool
}

func (s * Server) Start() {
	log.Printf("【Server】 listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)

	err := s.CallOnServerStartCallback()
	if err !=nil {
		log.Println("【Server】CallOnServerStartCallback err",err)
		s.stop<-true
		return
	}

	go func() {
		// start work pool
		s.MessageManager.StartWorkerPool()
		log.Printf("【Server】 worker pool starting")

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))

		if err != nil {
			log.Println("【Server】resolve tcp addr err: ", err)
			return
		}

		//get listener
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil{
			log.Println("【Server】listen", s.IPVersion, "err", err)
			return
		}

		log.Println("【Server】", s.Name ,"starting")

		var cid uint32
		cid = 0

		for  {
			// accept request
			accept, err := listener.AcceptTCP()
			if err != nil {
				log.Println("【Server】accept err", err)
				continue
			}

			//limit connection count
			if s.ConnectionManager.Count() > Conf["maxConn"].(int) {
				log.Println("【Server】connection limited")
				accept.Close()
				continue
			}

			conn := NewConnection(s, accept, cid, s.MessageManager, s.ConnectionManager)

			cid ++

			//one connection one goroutine
			go conn.Start()
		}
	}()
}

func (s *Server) Stop() {
	s.MessageManager.StopWorkerPool()
	s.ConnectionManager.Clean()
	s.CallOnServerStopCallback()
	s.stop <- true

	log.Println("【Server】has stopped")
}


func (s *Server) Run()  {

	s.Start()

	if <-s.stop {
		log.Println("【Server】has stopped")
		return
	}

}

func (s *Server) AddRouter(msgID uint32, router IRouter)  {

	s.MessageManager.AddRouter(msgID, router)
}

func (s *Server) SetOnServerStartCallback(callback func () error){
	s.OnServerStart = callback
}

func (s *Server) SetOnServerStopCallback(callback func () error){
	s.OnServerStop = callback
}

func (s *Server) CallOnServerStartCallback() error{
	if s.OnServerStart != nil {
		err := s.OnServerStart()
		return err
	}
	return nil
}

func (s *Server) CallOnServerStopCallback() error {
	if s.OnServerStop != nil {
		s.OnServerStop()
	}
	return nil
}


func (s *Server) SetOnConnStartCallback(callback func (IConnection)){
	s.OnConnStart = callback
}

func (s *Server) SetOnConnStopCallback(callback func (IConnection)){
	s.OnConnStop = callback
}

func (s *Server) CallOnConnStartCallback(conn IConnection){
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStopCallback(conn IConnection){
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}

func (s *Server) SetProperty(key string, value interface{})  {
	log.Println("【Server】set property", key)
	s.propertyLock.Lock()
	defer  s.propertyLock.Unlock()

	s.property[key] = value
}

func (s *Server) GetProperty(key string) (interface{} , error) {
	s.propertyLock.RLock()

	defer  s.propertyLock.RUnlock()

	if value, ok := s.property[key]; ok {
		return value, nil
	}else {
		return nil, errors.New("【Connection】property not found")
	}
}

func (s *Server) RemoveProperty(key string) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()

	delete(s.property, key)
}

func NewDefaultServer() *Server {
	return &Server{
		Name:    Conf["name"].(string),
		IP:      Conf["ip"].(string),
		Port:    Conf["port"].(int),
		IPVersion: Conf["ipVersion"].(string),
		MessageManager: NewMessageManager(),
		ConnectionManager:NewConnectionManager(),
		property:     make(map[string]interface{}),
		stop: make(chan bool,1),
	}
}

func NewServer(name string, host string, port int) *Server {
	return &Server{
		Name:    name,
		IP:      host,
		Port:    port,
		IPVersion: Conf["ipVersion"].(string),
		MessageManager: NewMessageManager(),
		ConnectionManager:NewConnectionManager(),
		property:     make(map[string]interface{}),
		stop: make(chan bool,1),
	}
}
