package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

func main() {

	grpcServer := grpc.NewServer()

	creds, err := credentials.NewServerTLSFromFile("src/grpc/tls/server.crt", "src/grpc/tls/server.key")
	if err != nil {
		log.Fatal(err)
	}
	//Adding filter
	grpcServer = grpc.NewServer(grpc.UnaryInterceptor(filter), grpc.Creds(creds))
	//grpcServer := grpc.NewServer(grpc.Creds(creds))
	//grpcServer := grpc.NewServer(

	RegisterPubsubServiceServer(grpcServer, NewPubsubService())

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
