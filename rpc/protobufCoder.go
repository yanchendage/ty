package rpc

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	pb2 "github.com/yanchendage/ty/demo/rpc/protobuf/pb"
	"github.com/yanchendage/ty/rpc/pb"
	"io"
	"reflect"
)

type ProtobufCoder struct {
	conn io.Closer
}

func (c *ProtobufCoder) DecodeHeader(buf []byte) (Header, error){
	//var msg Msg
	header := pb.RequestHeader{}

	// Marshal Header
	err := proto.Unmarshal(buf, &header)
	if err != nil {
		return Header{}, err
	}
	//这里需要拆分一下header和body
	buf = buf[0:len(buf)-int(header.Bodylen)]
	err = proto.Unmarshal(buf, &header)
	if err != nil {
		return Header{}, err
	}
	return Header{
		ServiceMethod:header.Method,
		Seq:header.Id,
	}, nil
}

func (c *ProtobufCoder) DecodeBody(buf []byte, bodyInterface interface{}) (reflect.Value, error) {
	//
	//var response proto.Message
	//if x != nil {
	//	var ok bool
	//	response, ok = x.(proto.Message)
	//	if !ok {
	//		return fmt.Errorf(
	//			"protorpc.ClientCodec.ReadResponseBody: %T does not implement proto.Message",
	//			x,
	//		)
	//	}
	//}
	var request proto.Message
	request, ok := bodyInterface.(proto.Message)

	if !ok {
		return reflect.Value{}, errors.New("bodyInterface does not implement proto.Message")
	}
	var pbRequest []byte

	if request != nil {
		err := proto.Unmarshal(pbRequest, request)
		if err != nil {
			return reflect.Value{}, err
		}
	}
	fmt.Println(1222)
	fmt.Println(pbRequest)
	fmt.Println(pbRequest)
	fmt.Println(request)

	return  reflect.ValueOf(request), nil
	//return  reflect.ValueOf(request).Elem(), nil

	//dec := gob.NewDecoder(bytes.NewBuffer(buf))
	//err := dec.Decode(&Header{})
	//err = dec.DecodeValue(reflect.ValueOf(argsInterface))
	//if err!=nil {
	//	return reflect.ValueOf(argsInterface), err
	//}
	//return reflect.ValueOf(argsInterface).Elem(), nil
	//
}


func (c *ProtobufCoder) DecodeResponse(buf []byte, bodyInterface interface{}) error{

	fmt.Println(buf)
	var request proto.Message
	request, ok := bodyInterface.(proto.Message)
	fmt.Println(request.String())

	if !ok {
		return  errors.New("bodyInterface does not implement proto.Message")
	}
	var pbRequest []byte

	if request != nil {
		err := proto.Unmarshal(pbRequest, request)
		if err != nil {
			return  err
		}
	}
	return nil
}
//func (c *ProtobufCoder) DecodeBody(buf []byte, bodyInterface interface{}) (reflect.Value, error) {
//
//	dec := gob.NewDecoder(bytes.NewBuffer(buf))
//	err := dec.Decode(&Header{})
//	err = dec.DecodeValue(reflect.ValueOf(argsInterface))
//	if err!=nil {
//		return reflect.ValueOf(argsInterface), err
//	}
//	return reflect.ValueOf(argsInterface).Elem(), nil
//}

//func (c *ProtobufCoder) DecodeReply(buf []byte, replyInterface interface{}) (reflect.Value, error) {
//
//	dec := gob.NewDecoder(bytes.NewBuffer(buf))
//	err := dec.Decode(&Header{})
//	err = dec.DecodeValue(reflect.ValueOf(replyInterface))
//	if err!=nil {
//		return reflect.ValueOf(replyInterface), err
//	}
//	return reflect.ValueOf(replyInterface).Elem(), nil
//}

func (c *ProtobufCoder) EncodeResponse(header *Header, body interface{}) ([]byte, error) {

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

	fmt.Println()

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



	fmt.Println(append(pbHeader, pbResponse...))

	h2 := &pb.RequestHeader{}
	proto.Unmarshal(append(pbHeader, pbResponse...),h2)

	fmt.Println(h2.Id)
	h3 := &pb2.SquareResponse{}
	proto.Unmarshal(pbResponse,h3)
	fmt.Println(pbResponse)

	fmt.Println(h3.Ans)





	return append(pbHeader, pbResponse...), nil




}



func (c *ProtobufCoder) EncodeRequest(header *Header, body interface{}) ([]byte, error) {

	//借鉴了protorpc
	var request proto.Message
	if body != nil {
		var ok bool
		if request, ok = body.(proto.Message); !ok {
			return nil, errors.New("body does not implement proto.Message")
		}
	}

	// marshal request
	pbRequest := []byte{}
	if request != nil {
		var err error
		pbRequest, err = proto.Marshal(request)
		if err != nil {
			return nil, err
		}
	}

	// generate header
	h := &pb.RequestHeader{
		Method:header.ServiceMethod,
		Id:header.Seq,
		Bodylen:uint64(len(pbRequest))}

	pbHeader, err := proto.Marshal(h)
	if err != err {
		return nil, err
	}

	//fmt.Println(pbHeader)
	//fmt.Println(pbRequest)
	//
	//h1 := &pb.RequestHeader{}
	//proto.Unmarshal(pbHeader,h1)
	//
	//fmt.Println("未拼接",h1.Id)
	//h2 := &pb.RequestHeader{}
	//pbMsg := append(pbHeader, pbRequest...)
	//
	//proto.Unmarshal(pbMsg,h2)
	//h3 := &pb2.SquareRequest{}
	//proto.Unmarshal(pbMsg,h3)
	//fmt.Println("dddd",h2.Method)
	//fmt.Println("ssss",h3.Num)
	//fmt.Println("sssszzzs",h2.Id)
	//fmt.Println("sssszzzs22",h2.Bodylen)
	//
	////分离pbHeader和pbRequest
	//
	//pbH := pbMsg[0:len(pbMsg)-int(h2.Bodylen)]
	//
	//h4 := &pb.RequestHeader{}
	//proto.Unmarshal(pbH,h4)
	//fmt.Println("h4",h4.Id)
	return append(pbHeader, pbRequest...), nil




}

//func (c *ProtobufCoder) Encode(msg Msg) ([]byte, error) {
//	buf := new(bytes.Buffer)
//	//proto.Marshal(pb.OrderRequest{})
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

//func (c *ProtobufCoder) ReadHeader(header *Header) error  {
//	return  c.decoder.Decode(header)
//}
//
//func (c *ProtobufCoder) ReadBody(body *Body) error  {
//	return  c.decoder.Decode(body)
//}
//
//func (c *ProtobufCoder) Write(header *Header, body *Body) error {
//	//defer func() {
//	//	//flush remain buffer to file
//	//	//err := c.buf.Flush()
//	//	if err!= nil {
//	//		c.Close()
//	//	}
//	//}()
//	//
//	//msg := Msg{
//	//	H: header,
//	//	B: body,
//	//}
//	//
//	//c.encoder.Encode(&msg)
//
//
//
//	//if err := c.encoder.Encode(header); err != nil{
//	//	log.Println("【RPC】gob encoding header err",err)
//	//	return err
//	//}
//	//
//	//if err := c.encoder.Encode(body); err != nil{
//	//	log.Println("【RPC】gob encoding body err",err)
//	//	return err
//	//}
//
//	return nil
//}

func (c *ProtobufCoder) Close() error {
	return c.conn.Close()
}

//Determine whether the interface is implemented
var _ ICoder = (*ProtobufCoder)(nil)

func NewProtobufCoder(conn io.Closer) ICoder {
	return &ProtobufCoder{
		conn:    conn,
	}
}



