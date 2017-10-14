package main

import (
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "proto"
)

const (
	address     = "localhost:50051"
)

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRedisClient(conn)

	r, err := c.Command(context.Background(), &pb.RedisRequest{Action:"get",Param:"2017"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Result)
}
