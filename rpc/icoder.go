package rpc

import (
	"io"
	"reflect"
)


//type Header struct {
//	ServiceMethod string // format "Service.Method"
//	Seq           uint64 // sequence number chosen by client
//	Error         string
//}

type ICoder interface {
	io.Closer

	//ReadHeader(header *Header) error
	//ReadBody(body *Body) error

	//Write(*Header, *Body) error

	EncodeRequest(*Header,interface{}) ([]byte, error)
	//DecodeRequestHeader(*Request) error
	//DecodeBody(interface{}) error
	//Encode(msg Msg) ([]byte, error)
	DecodeHeader(buf []byte) (Header, error)
	DecodeBody(buf []byte, body interface{}) (reflect.Value, error)

	EncodeResponse(*Header, interface{}) ([]byte, error)

	DecodeResponse(buf []byte, bodyInterface interface{}) error

	//DecodeArgs(buf []byte, argsInterface interface{}) (reflect.Value, error)
	//DecodeReply(buf []byte, replyInterface interface{}) (reflect.Value, error)
    //DecodeBody(buf []byte, interfaces ...interface{}) ([]reflect.Value, error)

	//Decode(msg *Msg)
}
