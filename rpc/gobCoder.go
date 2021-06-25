package rpc

import (
	"bytes"
	"encoding/gob"
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

func (c *GobCoder) Decode(buf []byte) (*Msg, error){
	var msg Msg
	dec := gob.NewDecoder(bytes.NewBuffer(buf))

	err := dec.Decode(&msg)

	if err != nil {
		return nil, err
	}

	return &msg ,err
}

func (c *GobCoder) Encode(msg Msg, register reflect.Type) ([]byte, error) {

	buf := new(bytes.Buffer)

	enc := gob.NewEncoder(buf)

	log.Println(register)

	gob.Register(reflect.New(register))

	err := enc.EncodeValue(msg.Args)
	//
	//if err == nil {
	//	fmt.Println(123)
	//}


	err := enc.Encode(&msg)

	if err != nil {
		return  nil, err
	}
	return buf.Bytes(), nil
}

func (c *GobCoder) ReadHeader(header *Header) error  {
	return  c.decoder.Decode(header)
}

func (c *GobCoder) ReadBody(body *Body) error  {
	return  c.decoder.Decode(body)
}

func (c *GobCoder) Write(header *Header, body *Body) error {
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
var _ ICoder = (*GobCoder)(nil)

func NewGobCoder(conn io.Closer) ICoder{

	var buf bytes.Buffer

	return &GobCoder{
		conn:    conn,
		//buf:     buf,
		encoder: gob.NewEncoder(&buf), // encode conn
		decoder: gob.NewDecoder(&buf), // decode conn
	}

}



