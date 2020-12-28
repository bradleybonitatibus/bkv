package server

import (
	"context"
	"fmt"
	"io"

	"github.com/bradleybonitatibus/bkv/db"
	pb "github.com/bradleybonitatibus/bkv/pb"
	"github.com/golang/protobuf/ptypes"
)

// Server object that contains the underlying gRPC server
type server struct {
	database db.Database
}

// Server is the gRPC server that will accept network requests to store and
// retrieve key value pairs from the database
type Server interface {
	Get(context.Context, *pb.GetRequest) (*pb.Response, error)
	Set(context.Context, *pb.SetRequest) (*pb.Response, error)
	BatchGet(pb.BKVService_BatchGetServer) error
	BatchSet(pb.BKVService_BatchSetServer) error
}

// NewServer initializes and returns a Server interface
func NewServer(db db.Database) Server {
	return &server{
		database: db,
	}
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	v, err := s.database.Get(req.GetKey())
	if err != nil {
		return nil, err
	}
	return &pb.Response{
		ResponseTimestamp: ptypes.TimestampNow(),
		Value:             v,
		Status:            pb.ResponseStatus_FOUND,
	}, nil
}

func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.Response, error) {
	k, v := req.GetKey(), req.GetValue()
	s.database.Set(k, v)
	return &pb.Response{
		ResponseTimestamp: ptypes.TimestampNow(),
		Value:             v,
		Status:            pb.ResponseStatus_CREATED,
	}, nil
}

// BatchGet method to implement bidirectional streaming of GET requests
func (s *server) BatchGet(stream pb.BKVService_BatchGetServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Println(err)
		}

		v, err := s.database.Get(req.GetKey())
		if err != nil {
			fmt.Println(err)
		}
		res := &pb.Response{
			ResponseTimestamp: ptypes.TimestampNow(),
			Value:             v,
			Status:            pb.ResponseStatus_FOUND,
		}
		if err := stream.Send(res); err != nil {
			fmt.Println(fmt.Sprintf("Failed sending message in batch get: %v", err))
		}
	}
}

func (s *server) BatchSet(stream pb.BKVService_BatchSetServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			fmt.Println(err)
		}

		k := req.GetKey()
		v := req.GetValue()

		s.database.Set(k, v)

		res := &pb.Response{
			ResponseTimestamp: ptypes.TimestampNow(),
			Value:             v,
			Status:            pb.ResponseStatus_CREATED,
		}
		if err := stream.Send(res); err != nil {
			fmt.Println("Failed sending response")
		}
	}
}
