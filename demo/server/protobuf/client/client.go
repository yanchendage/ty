package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yanchendage/ty/demo/server/protobuf/pb"
	"github.com/yanchendage/ty/server"
	"io"
	"log"
	"net"
	"time"
)

func main()  {
	conn,err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		log.Println("【Client】start err", err)
		return
	}
	coder := server.NewCoder()

	for {

		writeMsg, err := coder.Encode(server.NewMessage(0, []byte("fuck you ")))
		if err != nil {
			log.Println("【Client】encode msg err", err)
			return
		}

		_, err = conn.Write(writeMsg)
		if err != nil {
			log.Println("【Client】write msg err", err)
			return
		}

		log.Println("【Client】write msg success")

		head := make([]byte, coder.GetHeadLen())
		_, err = io.ReadFull(conn, head)

		if err != nil {
			log.Println("【Client】read head err", err)
			return
		}

		readMsgHead, err := coder.Decode(head)
		if err != nil {
			log.Println("【Client】head decode err", err)
			return
		}

		if readMsgHead.GetDataLen() > 0 {
			readMsg := readMsgHead.(*server.Message)
			readMsg.Data = make([]byte, readMsgHead.GetDataLen())

			_, err = io.ReadFull(conn,readMsg.Data)
			if err != nil{
				log.Println("【Client】read data err", err)
				return
			}

			fmt.Println("【Client】 Recv Msg: ID=", readMsg.ID, ", len=", readMsg.DataLen, ", data=", string(readMsg.Data))


			newdata := &pb.Pong{}
			err = proto.Unmarshal(readMsg.Data, newdata)
			if err != nil {
				fmt.Println("unmarshal err:", err)
			}

			fmt.Println(newdata.Msg)
		}

		time.Sleep(1*time.Second)
	}

}
