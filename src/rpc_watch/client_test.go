package rpc_watch

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"testing"
	"time"
)

func doClientWork(client *rpc.Client) {
	go func() {
		var keyChanged string
		err := client.Call("KVStoreService.Watch", 30, &keyChanged)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("watch:", keyChanged)
	}()

	time.Sleep(2 * time.Second)
	err := client.Call(
		"KVStoreService.Set", [2]string{"abcd", "abc-value"},
		new(struct{}),
	)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 3)
}

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	client := jsonrpc.NewClient(conn)
	doClientWork(client)
}
