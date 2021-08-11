package main

import "github.com/yanchendage/ty/rpc"

func main()  {
	var foo Foo

	serverManager := rpc.InitServerManagerAndRegister("RPC", "127.0.0.1", 7729,"http://127.0.0.1:8888/push")
	serverManager.RegisterService(&foo)
	serverManager.Run()
}