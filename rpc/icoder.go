package rpc

import (
	"io"
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
	DecodeRequestHeader(buf []byte) (Header, error)
	DecodeRequestBody(buf []byte, body interface{}) error

	EncodeResponse(*Header, interface{}) ([]byte, error)
	DecodeResponseHeader(buf []byte) (Header, error)
	DecodeResponseBody(buf []byte, bodyInterface interface{}) error

	//DecodeArgs(buf []byte, argsInterface interface{}) (reflect.Value, error)
	//DecodeReply(buf []byte, replyInterface interface{}) (reflect.Value, error)
    //DecodeBody(buf []byte, interfaces ...interface{}) ([]reflect.Value, error)

	//Decode(msg *Msg)
}
