package main

import (
	"fmt"
	"net"
	"os"

	"github.com/bradleybonitatibus/bkv/db"
	"github.com/bradleybonitatibus/bkv/pb"
	"github.com/bradleybonitatibus/bkv/server"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Failed to bind TCP listener")
		os.Exit(1)
	}

	s := grpc.NewServer()
	pb.RegisterBKVServiceServer(s, server.NewServer(db.NewDatabase()))
	if err := s.Serve(lis); err != nil {
		fmt.Println("Error occured while listening")
		os.Exit(1)
	}
}
