package rpc

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/yanchendage/ty/rpc/pb"
	"io"
	"reflect"
)

type GobCoder struct {
	conn io.Closer
	//buf bytes.Buffer
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func (c *GobCoder) DecodeHeader(buf []byte) (Header, error){
	//var msg Msg

	var header Header
	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	err := dec.Decode(&header)

	if err != nil {
		return Header{}, err
	}

	return header ,err
}

func (c *GobCoder) DecodeBody(buf []byte, interfaces interface{}) (reflect.Value, error) {

	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	err := dec.Decode(&Header{})
	err = dec.DecodeValue(reflect.ValueOf(interfaces))
	if err!=nil {
		return reflect.ValueOf(interfaces), err
	}
	return reflect.ValueOf(interfaces).Elem(), nil
}

//func (c *GobCoder) DecodeBody(buf []byte, interfaces ...interface{}) ([]reflect.Value, error) {
//
//	var reflectValueSlice []reflect.Value
//
//	dec := gob.NewDecoder(bytes.NewBuffer(buf))
//	err := dec.Decode(&Header{})
//
//	for _, i := range interfaces {
//		err = dec.DecodeValue(reflect.ValueOf(i))
//		if err!=nil {
//			return reflectValueSlice, err
//		}
//		reflectValueSlice = append(reflectValueSlice, reflect.ValueOf(i).Elem())
//	}
//
//	return reflectValueSlice, nil
//}

func (c *GobCoder) DecodeArgs(buf []byte, argsInterface interface{}) (reflect.Value, error) {

	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	err := dec.Decode(&Header{})
	err = dec.DecodeValue(reflect.ValueOf(argsInterface))
	if err!=nil {
		return reflect.ValueOf(argsInterface), err
	}
	return reflect.ValueOf(argsInterface).Elem(), nil
}

//func (c *GobCoder) DecodeReply(buf []byte, replyInterface interface{}) (reflect.Value, error) {
//
//	dec := gob.NewDecoder(bytes.NewBuffer(buf))
//	err := dec.Decode(&Header{})
//	err = dec.DecodeValue(reflect.ValueOf(replyInterface))
//	if err!=nil {
//		return reflect.ValueOf(replyInterface), err
//	}
//	return reflect.ValueOf(replyInterface).Elem(), nil
//}



func (c *GobCoder) EncodeRequest(header *Header, body interface{}) ([]byte, error) {

	buf := new(bytes.Buffer)

	enc := gob.NewEncoder(buf)

	err := enc.Encode(header)

	bodyV := reflect.ValueOf(body)

	if bodyV.IsValid() {
		err = enc.EncodeValue(bodyV)
	}
	//if msg.Reply.IsValid() {
	//	err = enc.EncodeValue(msg.Reply)
	//}
	if err != nil {
		return  nil, err
	}

	return buf.Bytes(), nil
}


//func (c *GobCoder) Encode(msg Msg) ([]byte, error) {
//
//	buf := new(bytes.Buffer)
//
//	enc := gob.NewEncoder(buf)
//
//	err := enc.Encode(msg.H)
//
//	if msg.Args.IsValid() {
//		err = enc.EncodeValue(msg.Args)
//	}
//	if msg.Reply.IsValid() {
//		err = enc.EncodeValue(msg.Reply)
//	}
//	if err != nil {
//		return  nil, err
//	}
//
//	return buf.Bytes(), nil
//}


//func (c *GobCoder) Write(header *Header, body *Body) error {
//
//	if err := c.encoder.Encode(header); err != nil{
//		log.Println("【RPC】gob encoding header err",err)
//		return err
//	}
//
//	if err := c.encoder.Encode(body); err != nil{
//		log.Println("【RPC】gob encoding body err",err)
//		return err
//	}
//
//	return nil
//}



func (c *GobCoder) EncodeResponse(header *Header, body interface{}) ([]byte, error) {

	//借鉴了protorpc
	var response proto.Message
	if body != nil {
		var ok bool
		if response, ok = body.(proto.Message); !ok {
			return nil, errors.New("body does not implement proto.Message")
		}
	}

	// marshal request
	pbResponse := []byte{}
	if response != nil {
		var err error
		pbResponse, err = proto.Marshal(response)
		if err != nil {
			return nil, err
		}
	}

	// generate header
	h := &pb.RequestHeader{
		Method:header.ServiceMethod,
		Id:header.Seq,
		Bodylen:uint64(len(pbResponse))}

	pbHeader, err := proto.Marshal(h)
	if err != err {
		return nil, err
	}

	//fmt.Println(pbHeader)
	//fmt.Println(pbRequest)
	//
	//h2 := &pb.RequestHeader{}
	//proto.Unmarshal(append(pbHeader, pbRequest...),h2)
	//h3 := &pb2.SquareRequest{}
	//proto.Unmarshal(append(pbHeader, pbRequest...),h3)
	//fmt.Println(h2.Method)
	//fmt.Println(h3.Num)




	return append(pbHeader, pbResponse...), nil




}


func (c *GobCoder) Close() error {
	return c.conn.Close()
}

////Determine whether the interface is implemented
//var _ ICoder = (*GobCoder)(nil)
//
//func NewGobCoder(conn io.Closer) ICoder {
//
//	var buf bytes.Buffer
//
//	return &GobCoder{
//		conn:    conn,
//		//buf:     buf,
//		encoder: gob.NewEncoder(&buf), // encode conn
//		decoder: gob.NewDecoder(&buf), // decode conn
//	}
//}
//


