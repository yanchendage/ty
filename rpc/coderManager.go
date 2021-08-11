package rpc

import (
	"io"
	"reflect"
)

//消息
type Msg struct {
	H Header
	Args reflect.Value
	Reply reflect.Value
}

type Mb struct {
	header []byte
	body []byte
}


//header
type Header struct {
	ServiceMethod string
	Seq uint64
	Err string
}

//body
type Body struct {
	Args reflect.Value
	Reply reflect.Value
}

// newCodeFunc is a func return Coder
type NewCodeFunc func(conn io.Closer) ICoder

type SerializationProtocol string

const GobProtocol SerializationProtocol  = "application/gob"
const JsonProtocol SerializationProtocol = "application/json"
const ProtobufProtocol SerializationProtocol = "application/protobuf"


var NewCoderFuncMap map[SerializationProtocol]NewCodeFunc

func init()  {
	NewCoderFuncMap = make(map[SerializationProtocol]NewCodeFunc)
	//NewCoderFuncMap[GobProtocol] = NewGobCoder
	NewCoderFuncMap[ProtobufProtocol] = NewProtobufCoder
}




