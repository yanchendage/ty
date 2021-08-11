package coder

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/yanchendage/ty/rpc"
	"io"
)

type JsonCoder struct {
	conn io.Closer
	buf bytes.Buffer
	encoder *json.Encoder
	decoder *json.Decoder
}

func (j *JsonCoder) Encode(msg rpc.Msg) ([]byte, error){
	return nil, nil
}

func (j *JsonCoder) Decode([]byte) (*rpc.Msg, error){
	return nil, nil
}

func (j *JsonCoder) ReadHeader(header *rpc.Header) error  {
	return  j.decoder.Decode(header)
}

func (j *JsonCoder) ReadBody(body *rpc.Body) error  {
	return  j.decoder.Decode(body)
}

func (j *JsonCoder) Write(header *rpc.Header, body *rpc.Body) error {
	return errors.New("123")
}

func (j *JsonCoder) Close() error {
	return j.conn.Close()
}

//Determine whether the interface is implemented
//var _ ICoder = (*JsonCoder)(nil)

func NewJsonCoder(conn io.Closer) rpc.ICoder {

	//var buf bytes.Buffer
	//
	//return &JsonCoder{
	//	//conn:    conn,
	//	////buf:     buf,
	//	//encoder: json.NewEncoder(&buf), // encode conn
	//	//decoder: json.NewDecoder(&buf), // decode conn
	//}

	return nil
}



