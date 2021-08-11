package main

import "github.com/yanchendage/ty/server"

func main()  {
	//new一个最简单的服务器
	r := server.NewServer("smiple","127.0.0.1", 7777)

	r.Run()
}
