package coder

import (
	"bytes"
	"encoding/gob"
	"github.com/yanchendage/ty/rpc"
	"io"
	"log"
	"reflect"
)

type GobCoder struct {
	conn io.Closer
	//buf bytes.Buffer
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func (c *GobCoder) DecodeHeader(buf []byte) (rpc.Header, error){
	//var msg Msg

	var header rpc.Header
	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	err := dec.Decode(&header)

	if err != nil {
		return rpc.Header{}, err
	}

	return header ,err
}

func (c *GobCoder) DecodeBody(buf []byte, interfaces ...interface{}) ([]reflect.Value, error) {

	var reflectValueSlice []reflect.Value

	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	err := dec.Decode(&rpc.Header{})

	for _, i := range interfaces {
		err = dec.DecodeValue(reflect.ValueOf(i))
		if err!=nil {
			return reflectValueSlice, err
		}
		reflectValueSlice = append(reflectValueSlice, reflect.ValueOf(i).Elem())
	}

	return reflectValueSlice, nil
}

func (c *GobCoder) DecodeArgs(buf []byte, argsInterface interface{}) (reflect.Value, error) {

	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	err := dec.Decode(&rpc.Header{})
	err = dec.DecodeValue(reflect.ValueOf(argsInterface))
	if err!=nil {
		return reflect.ValueOf(argsInterface), err
	}
	return reflect.ValueOf(argsInterface).Elem(), nil
}

func (c *GobCoder) DecodeReply(buf []byte, replyInterface interface{}) (reflect.Value, error) {

	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	err := dec.Decode(&rpc.Header{})
	err = dec.DecodeValue(reflect.ValueOf(replyInterface))
	if err!=nil {
		return reflect.ValueOf(replyInterface), err
	}
	return reflect.ValueOf(replyInterface).Elem(), nil
}

func (c *GobCoder) Encode(msg rpc.Msg) ([]byte, error) {

	buf := new(bytes.Buffer)

	enc := gob.NewEncoder(buf)

	err := enc.Encode(msg.H)

	if msg.Args.IsValid() {
		err = enc.EncodeValue(msg.Args)
	}
	if msg.Reply.IsValid() {
		err = enc.EncodeValue(msg.Reply)
	}
	if err != nil {
		return  nil, err
	}

	return buf.Bytes(), nil
}

func (c *GobCoder) ReadHeader(header *rpc.Header) error  {
	return  c.decoder.Decode(header)
}

func (c *GobCoder) ReadBody(body *rpc.Body) error  {
	return  c.decoder.Decode(body)
}

func (c *GobCoder) Write(header *rpc.Header, body *rpc.Body) error {
	//defer func() {
	//	//flush remain buffer to file
	//	//err := c.buf.Flush()
	//	if err!= nil {
	//		c.Close()
	//	}
	//}()
	//
	//msg := Msg{
	//	H: header,
	//	B: body,
	//}
	//
	//c.encoder.Encode(&msg)



	if err := c.encoder.Encode(header); err != nil{
		log.Println("【RPC】gob encoding header err",err)
		return err
	}

	if err := c.encoder.Encode(body); err != nil{
		log.Println("【RPC】gob encoding body err",err)
		return err
	}

	return nil
}

func (c *GobCoder) Close() error {
	return c.conn.Close()
}

//Determine whether the interface is implemented
var _ rpc.ICoder = (*GobCoder)(nil)

func NewGobCoder(conn io.Closer) rpc.ICoder {

	var buf bytes.Buffer

	return &GobCoder{
		conn:    conn,
		//buf:     buf,
		encoder: gob.NewEncoder(&buf), // encode conn
		decoder: gob.NewDecoder(&buf), // decode conn
	}
}



