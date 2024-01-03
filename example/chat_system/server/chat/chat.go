package chat

import (
	"golang.org/x/net/context"
	"grpc-test/protos"
	"log"
)

type Server struct {
	protos.UnimplementedChatServiceServer
}

func (s *Server) SayHello(ctx context.Context, msg *protos.Message) (*protos.Message, error) {
	log.Printf("Received message from client: [%s]", msg.Body)
	return &protos.Message{Body: "Hi from server!"}, nil
}

