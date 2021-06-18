package server

type Request struct {
	conn IConnection
	msg IMessage
	msgManager IMessageManager
}

func NewRequest(conn IConnection, msg IMessage, msgManage IMessageManager)  *Request{
	return &Request{
		conn:       conn,
		msg:        msg,
		msgManager: msgManage,
	}
}


func (r *Request) GetConnection() IConnection {
	return r.conn
}

func (r *Request) GetMsgData() []byte  {
	return  r.msg.GetData()
}

func (r *Request) GetMsgID() uint32  {
	return  r.msg.GetID()
}

//func (r *Request) task() {
//
//	var _ IWorker = (*Request)(nil)
//	r.msgManager.DoTask(r)
//
//}



