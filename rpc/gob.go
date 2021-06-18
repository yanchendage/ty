package rpc

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type Gob struct {
	conn 	io.ReadWriteCloser
	buf 	*bufio.Writer
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func (c *Gob) ReadHeader(header *Header) error  {
	return  c.decoder.Decode(header)
}

func (c *Gob) ReadBody(body interface{}) error  {
	return  c.decoder.Decode(body)
}

func (c *Gob) Write(header *Header, body interface{}) error {
	defer func() {
		//flush remain buffer to file
		err := c.buf.Flush()
		if err!= nil {
			c.Close()
		}
	}()

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

func (c *Gob) Close() error {
	return c.conn.Close()
}

//Determine whether the interface is implemented
var _ ICoder = (*Gob)(nil)

func NewGobCode(conn io.ReadWriteCloser) ICoder{
	buf := bufio.NewWriter(conn)

	return &Gob{
		conn:    conn,
		buf:     buf,
		encoder: gob.NewEncoder(buf), // encode buffer
		decoder: gob.NewDecoder(conn), // decode conn
	}
}



