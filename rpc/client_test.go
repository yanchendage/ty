package rpc

import (
	"log"
	"sync"
	"testing"
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

		log.Println("【RPC Client】id",i)

		wg.Add(1)

		go func(i int) {
			defer  wg.Done()
			args := Fuck{1,2}
			var reply string

			err := client.SyncCall("Foo.Sum", args, &reply)

			if err !=nil {
				log.Println("【Client】SyncCall err", err)
				return
			}

			log.Println("【Client】reply",reply)
		}(i)
	}

	wg.Wait()

	//conn,err := net.Dial("tcp", "127.0.0.1:7727")
	//if err != nil {
	//	log.Println("【Client】start err", err)
	//	return
	//}
	//
	//for i := 0; i < 5; i++{
	//	coder := server.NewCoder()
	//
	//	gobCoder := NewGobCoder(conn)
	//
	//	h := Header{
	//		ServiceMethod: "Fuck.You",
	//		Seq:          uint64(i),
	//	}
	//
	//	b := Body{
	//		Args:   []int{1,2,3},
	//	}
	//	msg := Msg{
	//		H: h,
	//		B: b,
	//	}
	//
	//	buf , err := gobCoder.Encode(msg)
	//
	//	if err !=nil {
	//		log.Println("【RPC】encode err ", err)
	//	}
	//
	//	//log.Println("【RPC】msg encode", buf)
	//	//
	//	//m, _ := gobCoder.Decode(buf)
	//
	//	//
	//	//log.Println("【RPC】msg decode", m)
	//
	//
	//
	//	writeMsg, err := coder.Encode(server.NewMessage(0, buf))
	//	if err != nil {
	//		log.Println("【Client】encode msg err", err)
	//		return
	//	}
	//
	//	_, err = conn.Write(writeMsg)
	//	if err != nil {
	//		log.Println("【Client】write msg err", err)
	//		return
	//	}
	//
	//	log.Println("【Client】write msg success")
	//}

}
