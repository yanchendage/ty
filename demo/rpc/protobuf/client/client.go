package main

import (
	"context"
	"github.com/yanchendage/ty/demo/rpc/protobuf/pb"
	"github.com/yanchendage/ty/rpc"
	"log"
)

func main()  {

	client, err := rpc.Dial("127.0.0.1:7729","application/protobuf")

	if err != nil {
		log.Println("【RPC Client】Dial err", err)
		return
	}

	defer  client.Close()

	args := &pb.SquareRequest{Num:12}
	reply := new(pb.SquareResponse)

	err = client.SyncCall(context.Background(), "Cal.Square", args, reply)
	if err !=nil {
		log.Println("【Client】SyncCall err", err)
		return
	}

	log.Println("arg",args,"reply",reply.Ans)
}