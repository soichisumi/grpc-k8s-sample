package api

import (
	"github.com/soichisumi/protobuf-trial/pbtrial"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":8080"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %+v\n", err)
	}
	server, _ := NewServer()
	s := grpc.NewServer()

	pbtrial.RegisterUserServiceServer(s, server)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}

}
