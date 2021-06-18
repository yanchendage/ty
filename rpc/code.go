package rpc

import "io"

type Header struct {
	ServerMethod string // server.method
	Seq uint64
	Err string
}

// newCodeFunc is a func retrun Coder
type NewCodeFunc func(closer io.ReadWriteCloser) ICoder

type Type string

const GobType Type  = "application/gob"
const JsonType Type = "application/json"


var NewCodeFuncMap map[Type]NewCodeFunc

func init()  {
	NewCodeFuncMap = make(map[Type]NewCodeFunc)
	NewCodeFuncMap[GobType] = NewGobCode
}




