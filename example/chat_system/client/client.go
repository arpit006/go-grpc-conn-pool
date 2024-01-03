package main

import (
	"context"
	v2 "github.com/arpit006/go-grpc-conn-pool/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	//"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials/insecure"
	"grpc-test/protos"
	"log"
	"time"
)

func main() {

	cfg := getGrpcClientConfig()

	// With pool
	conn, err := v2.NewClient(cfg, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Without pool
	//conn, err := grpc.Dial(":9003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Client could not connect to server on 9003. [%s]", err)
	}

	defer conn.Close()

	c := protos.NewChatServiceClient(conn)
	msg := &protos.Message{
		Body: "Server! Are you there??",
	}

	resp, err := c.SayHello(context.Background(), msg)
	if err != nil {
		log.Fatalf("error received from server. Error is: [%s]\n", err)
	}
	log.Printf("Response: [%s]", resp.Body)
}

func getGrpcClientConfig() *v2.ClientConfig {
	return v2.
		ClientConfigBuilder().
		WithName("grpc-test").
		WithTarget(":9003").
		WithPoolSize(3).
		WithConnMaxLifetime(2 * time.Minute).
		WithStdDeviation(10 * time.Second).
		Build()
}