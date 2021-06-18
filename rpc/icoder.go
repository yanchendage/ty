package rpc

import "io"

type ICoder interface {
	io.Closer

	ReadHeader(header *Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}
