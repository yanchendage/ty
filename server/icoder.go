package server

type ICoder interface {
	GetHeadLen()
	Encode()
	Decode()
}