package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/censture/unity-grpc-test-server/protobuf"
	"google.golang.org/grpc"
)

type helloServer struct {
}

func (s *helloServer) HelloStream(stream pb.UnityGRPC_HelloStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		go func(stream pb.UnityGRPC_HelloStreamServer, in *pb.HelloRequest) {
			i := 0
			for {
				i++

				res := &pb.HelloResponse{}
				res.Message = fmt.Sprintf("%d: Hello %s", i, in.GetName())
				err := stream.Send(res)
				if err != nil {
					log.Printf("error sending response: %v", err)
					break
				}

				time.Sleep(time.Second)
			}
		}(stream, req)
	}

	return nil
}

func (s *helloServer) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

func main() {
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUnityGRPCServer(s, &helloServer{})
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
