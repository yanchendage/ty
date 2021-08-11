package main

import (
	"context"
	"fmt"
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
	reply := &pb.SquareResponse{}


	err = client.SyncCall(context.Background(), "Cal.Square", args, reply)
	fmt.Println(123)
	if err !=nil {
		log.Println("【Client】SyncCall err", err)
		return
	}
	fmt.Println(reply)
	log.Println("arg",args,"reply",reply.Num)

}