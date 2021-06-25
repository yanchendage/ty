package server

type ICoder interface {
	GetHeadLen() uint32
	Encode(IMessage)([]byte, error)
	Decode(bin []byte) (IMessage, error)
}