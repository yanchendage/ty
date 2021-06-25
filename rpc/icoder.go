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

	Encode(msg Msg,register reflect.Type) ([]byte, error)
	Decode(buf []byte) (*Msg, error)

	//Decode(msg *Msg)
}
