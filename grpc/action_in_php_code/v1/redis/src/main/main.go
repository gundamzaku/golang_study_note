package main

import (
	"golang.org/x/net/context"
	pb "proto"
	"net"
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"fmt"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) Command(ctx context.Context, in *pb.RedisRequest) (*pb.RedisReply, error) {
	return &pb.RedisReply{Result: "Action is "+in.Action+"  Param is "+in.Param }, nil
}

func main() {
	fmt.Println("redis server start:");
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRedisServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
