package rpc

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestClient(t *testing.T) {

	client, err := Dial("127.0.0.1:7729","application/gob")

	if err != nil {
		log.Println("【RPC Client】Dial err", err)
		return
	}

	defer  client.Close()

	var wg sync.WaitGroup

	for i:=0; i<5 ;i++ {
		wg.Add(1)

		go func(i int) {
			defer  wg.Done()
			//args := &Param{1,2}
			arg2 := []string{"fuck","you"}
			//var reply int
			var reply2 string

			err := client.SyncCall(context.Background(), "Foo.Echo", arg2, &reply2)

			if err !=nil {
				log.Println("【Client】SyncCall err", err)
				return
			}

			log.Println("arg",arg2,"reply",reply2)

			//log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}

	wg.Wait()
}

func TestClientManager(t *testing.T)  {
	discovery := NewLocalDiscovery([]string{"127.0.0.1:7729"})
	cm := NewClientManager(discovery, Random)
	defer cm.Close()

	var wg sync.WaitGroup

	arg2 := []string{"fuck","you"}
	//var reply int
	var reply2 string

	for i :=1 ; i<5; i++{
		wg.Add(1)
		go func(i int) {
			defer  wg.Done()
			cm.Call(context.Background(),"Foo.Echo", arg2, &reply2 )

			fmt.Println(reply2)
		}(i)
	}

	wg.Wait()

}

func TestClientManagerRegistry(t *testing.T)  {
	discovery := NewRegistryDiscovery("http://127.0.0.1:8888", 10 * time.Second)
	cm := NewClientManager(discovery, Random)
	defer cm.Close()

	var wg sync.WaitGroup

	arg2 := []string{"fuck","you"}
	//var reply int
	var reply2 string

	for i :=1 ; i<5; i++{
		wg.Add(1)
		go func(i int) {
			defer  wg.Done()
			cm.Call(context.Background(),"Foo.Echo", arg2, &reply2 )

			fmt.Println(reply2)
		}(i)
	}

	wg.Wait()

}