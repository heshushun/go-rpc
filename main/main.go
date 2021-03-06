package main

import (
	"context"
	_ "encoding/json"
	"gorpc"
	_ "gorpc/codec"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func startServer(addr chan string) {
	var foo Foo
	if err := gorpc.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	gorpc.Accept(l)
}

func startServer2(addrCh chan string) {
	var foo Foo
	l, _ := net.Listen("tcp", ":9999")
	_ = gorpc.Register(&foo)
	gorpc.HandleHTTP()
	addrCh <- l.Addr().String()
	_ = http.Serve(l, nil)
}

func call(addrCh chan string) {
	client, _ := gorpc.DialHTTP("tcp", <-addrCh)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call(context.Background(), "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}

func main() {
	// day 1
	/*
		addr := make(chan string)
		go startServer(addr)

		// in fact, following code is like a simple gorpc client
		conn, _ := net.Dial("tcp", <-addr)
		defer func() { _ = conn.Close() }()

		time.Sleep(time.Second)
		// send options
		_ = json.NewEncoder(conn).Encode(gorpc.DefaultOption)
		cc := codec.NewGobCodec(conn)
		// send request & receive response
		for i := 0; i < 5; i++ {
			// send
			h := &codec.Header{
				ServiceMethod: "Foo.Sum",
				Seq:           uint64(i),
			}
			body := fmt.Sprintf("gorpc req %d", h.Seq)
			_ = cc.Write(h, body)
			_ = cc.ReadHeader(h)
			// receive
			var reply string
			_ = cc.ReadBody(&reply)
			log.Println("reply:", reply)
		}*/

	/*log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)

	client, _ := gorpc.Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		// day 2
		/*go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("gorpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Println("reply:", reply)
		}(i)
		// day 3
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()*/

	// day 5
	log.SetFlags(0)
	ch := make(chan string)
	go call(ch)
	startServer2(ch)
}
