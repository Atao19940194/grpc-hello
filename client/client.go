package main

import (
	"context"
	pb "grpc_demo/proto/helloworld"
	"log"

	"google.golang.org/grpc"
)

const PORT = "8080"

func main() {
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{
		NameTest: "",
	})

	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}

	log.Printf("resp: %s", resp.String())
}
