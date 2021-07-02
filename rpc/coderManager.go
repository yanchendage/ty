package rpc

import (
	"io"
	"reflect"
)

//var Byte bytes.Buffer

type Msg struct {
	H Header

	Args reflect.Value
	Reply reflect.Value

	//Args interface{}
	//Reply interface{}

	//B Body

	//MEthodType *methodType
	//svc *service
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


var NewCoderFuncMap map[SerializationProtocol]NewCodeFunc

func init()  {
	NewCoderFuncMap = make(map[SerializationProtocol]NewCodeFunc)
	NewCoderFuncMap[GobProtocol] = NewGobCoder
}




