package server

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Server struct {
	Name string
	IP string
	Port int
	IPVersion string
	MessageManager IMessageManager
	ConnectionManager IConnectionManager
}

func (s * Server) Start() {
	log.Printf("【Server】 listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)

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

			conn := NewConnection(accept, cid, s.MessageManager, s.ConnectionManager)

			cid ++

			//one connection one goroutine
			go conn.Start()
		}
	}()
}

func (s *Server) Stop() {

	log.Println("【Server】stop")
	s.ConnectionManager.Clean()
}


func (s *Server) Run()  {
	s.Start()

	for{
		time.Sleep(10 * time.Second)
	}

}

func (s *Server) AddRouter(msgID uint32, router IRouter)  {

	s.MessageManager.AddRouter(msgID, router)
}


func NewServer() *Server {
	return &Server{
		Name:    Conf["name"].(string),
		IP:      Conf["ip"].(string),
		Port:    Conf["port"].(int),
		IPVersion: Conf["ipVersion"].(string),
		MessageManager: NewMessageManager(),
		ConnectionManager:NewConnectionManager(),
	}
}