package server

type Message struct {
	ID uint32
	Data []byte
	DataLen uint32
}

func NewMessage(id uint32, data []byte ) *Message {
	return &Message{
		ID:      id,
		Data:    data,
		DataLen: uint32(len(data)),
	}
}

func (m *Message) GetID() uint32 {
	return m.ID
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}


func (m *Message) SetID(id uint32) {
	m.ID = id
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}