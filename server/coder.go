package server

import (
	"bytes"
	"encoding/binary"
)

type Coder struct {}

func NewCoder() *Coder  {
	return &Coder{}
}

func (c *Coder) GetHeadLen() uint32 {
	//id（4byte） + bodyLen（4byte）
	return 8
}

func (c *Coder) Encode(msg IMessage)([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	//write data length
	if err := binary.Write(buf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	//write id
	if err := binary.Write(buf, binary.LittleEndian, msg.GetID()); err != nil {
		return nil, err
	}
	//write data
	if err := binary.Write(buf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil ,err
	}

	return buf.Bytes(), nil
}

func (c *Coder) Decode(bin []byte) (IMessage, error) {
	reader := bytes.NewReader(bin)

	msg := &Message{}

	if err := binary.Read(reader, binary.LittleEndian, &msg.DataLen); err != nil{
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &msg.ID); err != nil{
		return nil, err
	}

	//todo valid msg

	return msg, nil

}