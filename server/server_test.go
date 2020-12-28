package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/bradleybonitatibus/bkv/db"
	"github.com/bradleybonitatibus/bkv/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	db := db.NewDatabase()
	pb.RegisterBKVServiceServer(s, NewServer(db))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestServer(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Error("Failed to dial bufnet")
	}
	defer conn.Close()

	client := pb.NewBKVServiceClient(conn)
	t.Run("Test Set Method", func(t *testing.T) {
		k, v := "Hello", "World"
		req := pb.SetRequest{
			Key:   k,
			Value: v,
		}
		res, err := client.Set(ctx, &req)
		if err != nil {
			t.Error("Failed to set value")
		}

		if res.GetStatus() != pb.ResponseStatus_CREATED {
			t.Errorf("Failed to have CREATED status in response")
		}

		if res.GetValue() != v {
			t.Error("Failed to set value")
		}
	})

	t.Run("Test GET key returns correct value", func(t *testing.T) {
		req := pb.GetRequest{
			Key: "Hello",
		}

		res, err := client.Get(ctx, &req)
		if err != nil {
			t.Error("Failed getting key")
		}
		if res.GetStatus() != pb.ResponseStatus_FOUND {
			t.Error("Failed to find key")
		}
		if res.GetValue() != "World" {
			t.Error("Failed to have correct value")
		}
	})

	t.Run("Test GET key returns error when key is not available", func(t *testing.T) {
		req := pb.GetRequest{
			Key: "Whomst",
		}

		res, err := client.Get(ctx, &req)
		if err == nil {
			t.Error("Expected an error")
		}
		fmt.Println(res)
	})
}
