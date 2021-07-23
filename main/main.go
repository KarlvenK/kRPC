package main

import (
	"fmt"
	kRPC "github.com/KarlvenK/krpc"
	"log"
	"net"
	"sync"
	"time"
)

func startServer(addr chan string) {
	//picl a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("newwork error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	kRPC.Accept(l)
}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)
	client, _ := kRPC.Dial("tcp", <-addr)
	defer func() {
		_ = client.Close()
	}()

	time.Sleep(time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("kRPC req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error: ", err)
			}
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
}
