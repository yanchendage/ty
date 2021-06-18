package server

type IMessage interface {
	GetID() uint32
	GetData() []byte
	GetDataLen() uint32

	SetID(uint32)
	SetData([]byte)
	SetDataLen(uint32)
}
