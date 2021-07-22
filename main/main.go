package main

import (
	"encoding/json"
	"fmt"
	kRPC "github.com/KarlvenK/krpc"
	"github.com/KarlvenK/krpc/codec"
	"log"
	"net"
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
	addr := make(chan string)
	go startServer(addr)

	conn, _ := net.Dial("tcp", <-addr)
	defer func() {
		_ = conn.Close()
	}()

	time.Sleep(time.Second)

	_ = json.NewEncoder(conn).Encode(kRPC.DefaultOption)
	cc := codec.NewGobCodec(conn)

	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("krpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}
