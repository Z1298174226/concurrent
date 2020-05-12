package rpc_watch

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"testing"
)

func TestServer(t *testing.T) {

	KVStore := NewKVStoreService()
	rpc.Register(KVStore)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
		}
		go jsonrpc.ServeConn(conn)
	}
}
