package main

import "github.com/yanchendage/ty/rpc"

func main()  {
	rpc.NewLocalRegistryServer("127.0.0.1:8888")
}
