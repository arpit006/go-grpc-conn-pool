package main

import (
	"grpc-test/protos"
	"grpc-test/server/chat"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Trying to start server.......")
	lis, err := net.Listen("tcp", ":9003")
	if err != nil {
		log.Fatalf("Failed to listen on port 9003: %v", err)
	}

	grpcServer := grpc.NewServer()

	s := chat.Server{}

	protos.RegisterChatServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to server gRPC server over port 9003: %v", err)
	}
	log.Println("Starting gRPC server on port 9003")
}
