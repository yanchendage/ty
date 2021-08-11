package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yanchendage/ty/demo/server/protobuf/pb"
	"github.com/yanchendage/ty/server"
	"log"
)


type PongRouter struct {
	server.BRouter
}

func (this *PongRouter) Handle(request server.IRequest) {
	log.Println("【Request】 msgId=", request.GetMsgID(), ", data=", string(request.GetMsgData()))

	pong := &pb.Pong{Msg:"pong"}

	//person := &pb.Person{
	//	Name:   "Aceld",
	//	Age:    16,
	//	Emails: []string{"https://legacy.gitbook.com/@aceld", "https://github.com/aceld"},
	//	Phones: []*pb.PhoneNumber{
	//		&pb.PhoneNumber{
	//			Number: "13113111311",
	//			Type:   pb.PhoneType_MOBILE,
	//		},
	//		&pb.PhoneNumber{
	//			Number: "14141444144",
	//			Type:   pb.PhoneType_HOME,
	//		},
	//		&pb.PhoneNumber{
	//			Number: "19191919191",
	//			Type:   pb.PhoneType_WORK,
	//		},
	//	},
	//}


	data, err := proto.Marshal(pong)
	if err != nil {
		fmt.Println("marshal err:", err)
	}

	err = request.GetConnection().SendMsg(0, data)
	if err != nil {
		fmt.Println(err)
	}
}
func main()  {
	//new一个最简单的服务器
	s := server.NewServer("protobuf","127.0.0.1", 7777)
	s.AddRouter(0, &PongRouter{})
	s.Run()
}
