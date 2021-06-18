package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type Connection struct {
	ConnectionManager IConnectionManager

	Conn *net.TCPConn

	ID uint32

	Closed bool

	Exit chan bool

	MessageManager IMessageManager

	Msg chan []byte

	BuffMsg chan []byte
}

func NewConnection(conn *net.TCPConn, id uint32, messageManager IMessageManager, connectionManager IConnectionManager) *Connection {
	connection := &Connection{
		ConnectionManager:connectionManager,
		Conn:   conn,
		ID:     id,
		Closed: false,
		MessageManager : messageManager,
		Exit:   make(chan bool, 1),
		Msg: make(chan []byte),
		BuffMsg : make(chan []byte, Conf["buffMessageQueueSize"].(int)),
	}

	//add connection to connection manager
	connectionManager.Add(connection)

	return connection
}

func (c *Connection) read() {
	log.Println("【Connection】",c.GetRemoteAddr().String()," start read")
	defer log.Println("【Connection】",c.GetRemoteAddr().String()," end read")
	defer  c.Stop()

	coder := NewCoder()

	for {
		headData := make([]byte, coder.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, headData); err != nil {
			fmt.Println("read msg head error ", err)
			return
		}

		msg, err := coder.Decode(headData)
		if err != nil {
			log.Println("【Connection】coder decode err ", err)
			c.Exit <- true
			return
		}

		var body []byte
		if msg.GetDataLen() > 0 {
			body = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConn(),body); err != nil{
				log.Println("【Connection】read msg body err ", err)
				c.Exit <- true
				return
			}

			msg.SetData(body)
		}

		req := NewRequest(c, msg, c.MessageManager)

		var _ IRequest = (*Request)(nil)

		c.MessageManager.AddWorkerToWorkQueue(req)
	}
}

func (c *Connection) write()  {
	for {
		select {
		case data := <-c.Msg:

			if _, err := c.Conn.Write(data); err != nil {
				log.Println("【Connection】send msg err ", err)
				return
			}

			log.Println("【Connection】send msg success")
		case data, ok := <- c.BuffMsg:

			if ok {
				_, err := c.Conn.Write(data)
				if err != nil {
					log.Println("【Connection】send buff msg err ", err)
					return
				}

				log.Println("【Connection】send buff msg success")
			}else {
				log.Println("【Connection】send buff msg success")
			}
			
		case <- c.Exit:
			return
		}
	}
}

func (c *Connection) SendMsg(msgID uint32, data []byte) error{
	if c.Closed == true {
		return errors.New("【Connection】closed when call SendMsg")
	}

	coder := NewCoder()

	msg ,err := coder.Encode(NewMessage(msgID, data))
	if err != nil {
		log.Println("【Connection】msg encode err", err)
		return  errors.New("【Connection】msg encode err when send msg")
	}

	c.Msg <- msg

	return  nil
}

func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error{
	if c.Closed == true {
		return errors.New("【Connection】closed when call SendBuffMsg")
	}

	coder := NewCoder()

	msg ,err := coder.Encode(NewMessage(msgID, data))
	if err != nil {
		log.Println("【Connection】buffMsg encode err", err)
		return  errors.New("【Connection】buffMsg encode err when send msg")
	}

	c.BuffMsg <- msg

	return  nil
}


func (c *Connection) Start()  {
	go c.read()
	go c.write()

	for {
		select {
		case <- c.Exit:
			return
		}
	}

}

func (c *Connection) Stop() {
	log.Println("【Connection】id=", c.ID, "closing")
	if c.Closed == true {
		return
	}

	c.Closed = true

	c.Conn.Close()

	c.Exit <- true

	c.ConnectionManager.Remove(c.ID)

	close(c.Exit)
	close(c.BuffMsg)
}

func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ID
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}


