package main

import (
	"context"
	"fmt"
	"github.com/yanchendage/ty/rpc"
	"sync"
	"time"
)

func main()  {
	discovery := rpc.NewRegistryDiscovery("http://127.0.0.1:8888", 10 * time.Second)
	cm := rpc.NewClientManager(discovery, rpc.Random)
	defer cm.Close()

	var wg sync.WaitGroup

	arg2 := []string{"fuck","you"}
	//var reply int
	var reply2 string

	for i := 0; i< 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer  wg.Done()
			cm.Call(context.Background(),"Foo.Echo", arg2, &reply2 )

			fmt.Println(reply2)
		}(i)
	}

	wg.Wait()
}