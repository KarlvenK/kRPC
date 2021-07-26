package main

import (
	"context"
	kRPC "github.com/KarlvenK/krpc"
	"log"
	"net"
	"sync"
	"time"
)

/*
//day 1
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
*/

//day 3

type Foo int
type Args struct {
	Num1, Num2 int
}

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num2 + args.Num1
	return nil
}

func startService(addr chan string) {
	var foo Foo
	if err := kRPC.Register(&foo); err != nil {
		log.Fatal("register error: ", err)
	}

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error: ", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	kRPC.Accept(l)
}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startService(addr)
	client, _ := kRPC.Dial("tcp", <-addr)
	defer func() {
		_ = client.Close()
	}()

	time.Sleep(time.Second)
	//send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{
				Num1: i,
				Num2: i * i,
			}
			var reply int
			if err := client.Call(context.Background(), "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}
