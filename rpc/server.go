package rpc

import (
	"encoding/json"
	"github.com/yanchendage/ty/server"
	"io"
	"log"
	"net"
	"reflect"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int
	Type Type
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	Type: GobType,
}

func NewS()  {
	r := server.NewServer()
	r.AddRouter()

	r.Run()
}

type Server struct {
}

func NewServer() *Server  {
	return &Server{}
}

var DefaultServer = NewServer()

func (server * Server) Accept(lis net.Listener)  {

	for  {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("【RPC】: error accept:", err)
			return
		}
		go server.ReadAndWrite(conn)
	}
}

func(server * Server) ReadAndWrite(conn io.ReadWriteCloser)  {
	defer func() {
		conn.Close()
	}()

	var opt Option

	if err := json.NewDecoder(conn).Decode(&opt); err != nil{
		log.Println("【RPC】: invalid Sign:", err)
		return
	}

	f := NewCodeFuncMap[opt.Type]
	if f == nil {
		if err := json.NewDecoder(conn).Decode(&opt); err != nil{
			log.Println("【RPC】: invalid Type:", err)
			return
		}
	}

}


type request struct {
	header 	*Header
	arg    	reflect.Value
	reply	reflect.Value
}

func (server *Server) readRequestHeader(coder ICoder)(*Header, error) {
	var h Header
	if err := coder.ReadHeader(&h); err!=nil{
		if err != io.EOF && err != io.ErrUnexpectedEOF{
			log.Println("【RPC】: error read header :", err)
		}
		return nil, err
	}

	return &h, nil
}

func (server Server) readRequest(coder ICoder)  (*request, error){
	h, err := server.readRequestHeader(coder)
	if err != nil {
		return nil, err
	}

	req := &request{header:h}

	req.arg = reflect.New(reflect.TypeOf(""))
	err = coder.ReadBody(req.arg.Interface())
	if err != nil {
		log.Println("【RPC】: error read body :", err)
	}
	return req, nil
}
