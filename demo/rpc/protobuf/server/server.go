package main

import (
	"github.com/yanchendage/ty/demo/rpc/protobuf/pb"
	"github.com/yanchendage/ty/rpc"
)

type Cal int

func (cal *Cal) Square(request *pb.SquareRequest,response *pb.SquareResponse) error {
	response.Num = request.Num
	response.Ans = request.Num * request.Num
	return nil
}

func main()  {
	var cal Cal

	serverManager := rpc.InitServerManager("RPC", "127.0.0.1", 7729)
	serverManager.RegisterService(&cal)

	serverManager.Run()
}
