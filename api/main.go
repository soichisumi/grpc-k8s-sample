package main

import (
	"fmt"
	"github.com/soichisumi/grpc-auth-sample/api-pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":3000"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %+v\n", err)
	}
	server, err := NewServer("./privKey.pem", "./privKey.pem.pub.pkcs8")
	if err != nil{
		log.Fatalf("failed to create server: %+v\n", err)
	}
	s := grpc.NewServer()

	apipb.RegisterUserServiceServer(s, server)
	reflection.Register(s)
	fmt.Printf("grpc server is running on port:%s...\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}

}
