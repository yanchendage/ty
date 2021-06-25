package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"reflect"
)

type JsonCoder struct {
	conn io.Closer
	buf bytes.Buffer
	encoder *json.Encoder
	decoder *json.Decoder
}

func (j *JsonCoder) Encode(msg Msg, register reflect.Type) ([]byte, error){
	return nil, nil
}

func (j *JsonCoder) Decode([]byte) (*Msg, error){
	return nil, nil
}

func (j *JsonCoder) ReadHeader(header *Header) error  {
	return  j.decoder.Decode(header)
}

func (j *JsonCoder) ReadBody(body *Body) error  {
	return  j.decoder.Decode(body)
}

func (j *JsonCoder) Write(header *Header, body *Body) error {
	return errors.New("123")
}

func (j *JsonCoder) Close() error {
	return j.conn.Close()
}

//Determine whether the interface is implemented
//var _ ICoder = (*JsonCoder)(nil)

func NewJsonCoder(conn io.Closer) ICoder{

	var buf bytes.Buffer

	return &JsonCoder{
		conn:    conn,
		//buf:     buf,
		encoder: json.NewEncoder(&buf), // encode conn
		decoder: json.NewDecoder(&buf), // decode conn
	}

}



