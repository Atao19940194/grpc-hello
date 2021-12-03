package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"

	helloworldpb "grpc_demo/proto/helloworld"
	"grpc_demo/route"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type server struct {
	helloworldpb.UnimplementedGreeterServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) SayHello(ctx context.Context, in *helloworldpb.HelloRequest) (*helloworldpb.HelloReply, error) {
	// err := protojson.Unmarshal(bytes, in)
	return &helloworldpb.HelloReply{TestMessage: in.NameTest + " world"}, nil
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	helloworldpb.RegisterGreeterServer(s, &server{})
	// Serve gRPC Server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"127.0.0.1:8080",
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)

	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	/*----------------------------------------------------------------------------------------------------*/
	port := 8090

	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}))

	// Register Greeter
	err = helloworldpb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	router := gin.Default()
	router.Any("/rpc/v1/*any", gin.Logger(), gin.Recovery(), gin.WrapF(gwmux.ServeHTTP))
	route.InitRoute(router)
	err = endless.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
