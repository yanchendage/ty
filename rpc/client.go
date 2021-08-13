package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/yanchendage/ty/server"
	"io"
	"log"
	"net"
	"sync"
)

type Caller struct {
	Seq uint64
	ServiceMethod string
	Args interface{}
	Reply interface{}
	Error error
	Done chan *Caller
}

func (c *Caller) done()  {
	c.Done <- c
}

type Client struct {
	seq uint64
	conn *net.TCPConn
	Coder ICoder
	ServerCoder server.ICoder
	mu sync.Mutex
	sending sync.Mutex
	pending map[uint64]*Caller
	closed bool
}

var _ io.Closer = (*Client)(nil)

func (client *Client) Close() error  {
	client.mu.Lock()
	defer  client.mu.Unlock()

	if client.closed {
		return errors.New("【RPC client】 client is closed")
	}
	client.closed = true

	return client.conn.Close()
}

func (client *Client) AddCaller(caller *Caller) (uint64, error) {
	client.mu.Lock()
	defer  client.mu.Unlock()

	if client.closed {
		return 0, errors.New("【RPC client】 client is closed")
	}
	caller.Seq = client.seq
	client.pending[caller.Seq] = caller
	client.seq ++ ;

	return  caller.Seq, nil
}


func (client *Client) RemoveCaller(seq uint64) *Caller {
	client.mu.Lock()
	defer client.mu.Unlock()

	caller := client.pending[seq]
	delete(client.pending, seq)

	return caller
}

func (client *Client) CleanCaller(err error) {
	client.mu.Lock()
	defer  client.mu.Unlock()

	client.closed = true

	for _, caller := range client.pending {
		caller.Error = err
		caller.done()
	}
}

func (client *Client) receive()  {

	for {
		head := make([]byte, client.ServerCoder.GetHeadLen())
		_, err := io.ReadFull(client.conn, head)

		if err != nil {
			log.Println("【RPC Client】read head err", err)
			return
		}

		readMsgHead, err := client.ServerCoder.Decode(head)
		if err != nil {
			log.Println("【RPC Client】head decode err", err)
			return
		}
		if readMsgHead.GetDataLen() > 0 {
			readMsg := readMsgHead.(*server.Message)

			readMsg.Data = make([]byte, readMsgHead.GetDataLen())

			_, err = io.ReadFull(client.conn,readMsg.Data)
			if err != nil{
				log.Println("【RPC Client】read data err", err)
				return
			}
			header, err := client.Coder.DecodeResponseHeader(readMsg.Data)

			if err != nil {
				log.Println("【RPC Client】data decode header err", err)
				return
			}

			caller := client.RemoveCaller(header.Seq)

			switch {
			case caller == nil:
				fmt.Println("caller == nil")
				err = errors.New("【RPC Client】caller was already removed")
			case header.Err != "":
				fmt.Println(header.Err != "")
				caller.Error = fmt.Errorf(header.Err)
				caller.done()
			default:
				//var replyv reflect.Value
				//replyType := reflect.TypeOf(caller.Reply)
				//if replyType.Kind() == reflect.Ptr {
				//	replyv = reflect.New(replyType.Elem())
				//} else {
				//	replyv = reflect.New(replyType).Elem()
				//}
				//
				//replyvi := replyv.Interface()
				//
				//if replyv.Type().Kind() != reflect.Ptr {
				//	replyvi = replyv.Addr().Interface()
				//}
				//
				//fmt.Println(replyvi)
				//
				//reply, err := client.Coder.DecodeBody(readMsg.GetData(),replyvi)
				////reply,err := client.Coder.DecodeReply(readMsg.GetData(),argvi)
				//if err != nil{
				//	log.Println("【RPC Client】data decode body err", err)
				//	return
				//}
				////
				//reflect.ValueOf(caller.Reply).Elem().Set(reply)
				//
				//caller.done()

				//client.Coder.DecodeRe

				//fmt.Println("readMsg.GetData()",readMsg.GetData())
				//reply, err := client.Coder.DecodeBody(readMsg.GetData(),caller.Reply)
				//if err != nil{
				//	log.Println("【RPC Client】data decode body err", err)
				//	return
				//}
				//fmt.Println("sds")
				//caller.done()
				//fmt.Println(reply)

				err := client.Coder.DecodeResponseBody(readMsg.GetData(),caller.Reply)
				if err != nil{
					log.Println("【RPC Client】data decode body err", err)
					return
				}
				//fmt.Println("sds")
				caller.done()
				//for i := 0; i < reply.Elem().NumField(); i++ {
				//
				//	field := typ.Elem().Field(i)  //字段的数据类型
				//	value := val.Elem().Field(i)  //字段的数据值
				//	fmt.Println("type1:", field)  //type: {Name  string json:"name" 0 [0] false}
				//	fmt.Println("value1:", value) //value: lei
				//	switch value.Kind() {
				//	case reflect.Int:
				//		value.SetInt(88) //往该字段设值
				//	case reflect.String:
				//		value.SetString("Test") // 往该字段设值
				//	default:
				//		fmt.Println("类型不支持")
				//	}
				//	fmt.Println("type2:", field)  //type: {Name  string json:"name" 0 [0] false}
				//	fmt.Println("value2:", value) //value: Test
				//	fmt.Println(field.Tag.Get("json"))
				//}
				//
				//
				//caller.Reply
				//fmt.Println(caller.Reply)
				//fmt.Println(reply)
				//reflect.ValueOf(caller.Reply)..Elem().Set(reply)
				//reflect.ValueOf(caller.Reply).Elem().Set(reply)
			}
			//fmt.Println("【RPC Client】 Recv Msg: ID=", readMsg.ID, ", len=", readMsg.DataLen, ", data=", string(readMsg.Data))
		}
	}
}

func Dial(address string, serializationProtocol SerializationProtocol) (*Client, error) {

	conn, err := net.Dial("tcp", address)

	if err != nil {
		log.Println("【RPC Client】net Dial err", err)
		return nil, err
	}

	f := NewCoderFuncMap[serializationProtocol]

	if f == nil {
		log.Println("【RPC Client】 invilad serialization protocol", err)
		return nil, err
	}

	serverCoder := server.NewCoder()
	rpcCoder := f(conn)

	client := &Client{
		seq : 1,
		Coder : rpcCoder,
		ServerCoder: serverCoder,
		pending: make(map[uint64]*Caller),
		conn : conn.(*net.TCPConn),
	}

	//read and write

	go client.receive()

	return client, err
}

//第二版
func (client *Client) Send(caller *Caller) {

	client.sending.Lock()
	defer client.sending.Unlock()

	seq , err := client.AddCaller(caller)

	if err != nil {
		log.Println("【RPC Client】AddCaller err", err)
		caller.Error = err
		caller.done()
		return
	}

	h := Header{
		ServiceMethod: caller.ServiceMethod,
		Seq:          seq,
	}

	b, err := client.Coder.EncodeRequest(&h,caller.Args)

	if err != nil {
		log.Println("【RPC Client】Encode err", err)
		return
	}

	writeMsg, err := client.ServerCoder.Encode(server.NewMessage(1, b))

	if err != nil {
		log.Println("【RPC Client】encode msg err", err)
		return
	}

	_, err = client.conn.Write(writeMsg)
	if err != nil {
		log.Println("【RPC Client】write msg err", err)
		return
	}

	if err != nil {
		log.Println("【RPC Client】send msg err", err)
		return
	}
	log.Println("【RPC Client】send msg success", b)
}



//第一版代码
//func (client *Client) Send(caller *Caller) {
//
//	client.sending.Lock()
//	defer client.sending.Unlock()
//
//	seq , err := client.AddCaller(caller)
//
//	if err != nil {
//		log.Println("【RPC Client】AddCaller err", err)
//		caller.Error = err
//		caller.done()
//		return
//	}
//
//	h := Header{
//		ServiceMethod: caller.ServiceMethod,
//		Seq:          seq,
//	}
//
//	msg := Msg{
//		H: h,
//		Args:reflect.ValueOf(caller.Args),
//	}
//
//	buf , err := client.Coder.Encode(msg)
//
//	if err !=nil {
//		log.Println("【RPC】encode err ", err)
//	}
//
//		writeMsg, err := client.ServerCoder.Encode(server.NewMessage(0, buf))
//		if err != nil {
//			log.Println("【Client】encode msg err", err)
//			return
//		}
//
//		_, err = client.conn.Write(writeMsg)
//		if err != nil {
//			log.Println("【Client】write msg err", err)
//			return
//		}
//
//	if err != nil {
//		log.Println("【RPC Client】send msg err", err)
//		return
//	}
//	log.Println("【RPC Client】send msg success", buf)
//}

func (client *Client) AsyncCall(serviceMethod string, args interface{}, reply interface{}, done chan *Caller) *Caller{
	if done == nil{
		done = make(chan *Caller, 1)
	} else if cap(done) == 0 {
		log.Panic("【RPC client】 done channel is unbuffered")
	}

	caller := &Caller{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}

	client.Send(caller)
	return caller
}

func (client *Client)Available() bool  {
	client.mu.Lock()
	defer client.mu.Unlock()

	return !client.closed
}

func (client *Client) SyncCall(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error{

	call := client.AsyncCall(serviceMethod, args, reply, make(chan *Caller, 1))

	for  {
		select {
		case <-ctx.Done():
			client.RemoveCaller(call.Seq)
			return errors.New("rpc client: call failed: " + ctx.Err().Error())
		case call := <-call.Done:
			return call.Error
		}
	}
}



