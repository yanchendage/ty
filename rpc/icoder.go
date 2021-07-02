package rpc

import (
	"io"
	"reflect"
)

type ICoder interface {
	io.Closer

	ReadHeader(header *Header) error
	ReadBody(body *Body) error

	Write(*Header, *Body) error

	Encode(msg Msg) ([]byte, error)
	DecodeHeader(buf []byte) (Header, error)
	DecodeArgs(buf []byte, argsInterface interface{}) (reflect.Value, error)
	DecodeReply(buf []byte, replyInterface interface{}) (reflect.Value, error)

    DecodeBody(buf []byte, interfaces ...interface{}) ([]reflect.Value, error)

	//Decode(msg *Msg)
}
