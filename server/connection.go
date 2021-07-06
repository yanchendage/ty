package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Connection struct {
	Server IServer
	ConnectionManager IConnectionManager
	Conn *net.TCPConn
	ID uint32
	Closed bool
	MessageManager IMessageManager
	Msg chan []byte
	BuffMsg chan []byte
	ctx  context.Context
	cancel context.CancelFunc
	property map[string]interface{}
	propertyLock sync.RWMutex
}

func NewConnection(server IServer, conn *net.TCPConn, id uint32, messageManager IMessageManager, connectionManager IConnectionManager) *Connection {
	connection := &Connection{
		Server: server,
		ConnectionManager:connectionManager,
		Conn:   conn,
		ID:     id,
		Closed: false,
		MessageManager : messageManager,
		Msg: make(chan []byte),
		BuffMsg : make(chan []byte, Conf["buffMessageQueueSize"].(int)),
		property:     make(map[string]interface{}), 
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
		select {
		case <-c.ctx.Done():
			return
		default:

			headData := make([]byte, coder.GetHeadLen())
			if _, err := io.ReadFull(c.Conn, headData); err != nil {
				//if err != io.EOF && err != io.ErrUnexpectedEOF {
				//	fmt.Println("【Connection】read msg head error ", err)
				//	return
				//}

				if err != io.EOF {
					fmt.Println("【Connection】read msg head error ", err)
					return
				}

				return
			}

			msg, err := coder.Decode(headData)
			if err != nil {
				log.Println("【Connection】coder decode err ", err)
				return
			}

			var body []byte
			if msg.GetDataLen() > 0 {
				body = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.GetTCPConn(),body); err != nil{
					log.Println("【Connection】read msg body err ", err)
					return
				}

				msg.SetData(body)
			}

			req := NewRequest(c, msg, c.MessageManager)

			var _ IRequest = (*Request)(nil)

			c.MessageManager.AddWorkerToWorkQueue(req)

		}
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
		case <-c.ctx.Done():
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

	c.ctx, c.cancel = context.WithCancel(context.Background())

	go c.read()
	go c.write()

	c.Server.CallOnConnStartCallback(c)
}

func (c *Connection) Stop() {
	log.Println("【Connection】id=", c.ID, "closing")
	if c.Closed == true {
		return
	}
	c.Server.CallOnConnStopCallback(c)

	c.Closed = true

	c.Conn.Close()

	c.cancel()

	c.ConnectionManager.Remove(c.ID)

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

func (c * Connection) SetProperty(key string, value interface{})  {
	c.propertyLock.Lock()
	defer  c.propertyLock.Unlock()

	c.property[key] = value
}

func (c * Connection) GetProperty(key string) (interface{} , error) {
	c.propertyLock.RLock()

	defer  c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}else {
		return nil, errors.New("【Connection】property not found")
	}
}

func (c * Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func (c *Connection) GetServer() IServer {
	return c.Server
}


