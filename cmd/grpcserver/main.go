package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"

	"google.golang.org/grpc"

	"grpc-test/constants"
	api "grpc-test/gen/proto/go/results/api/v1"
	"grpc-test/s3"
)

type Server struct {
	api.UnimplementedTestServiceServer
}

func SetupGrpcServer() *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.MaxSendMsgSize(math.MaxInt32),
	}
	grpcServer := grpc.NewServer(opts...)
	api.RegisterTestServiceServer(grpcServer, &Server{})

	return grpcServer
}

func getS3Reader(bucket string, key string) io.Reader {
	ctx := context.Background()
	manager := s3.NewS3Manager(ctx)
	return manager.GetS3Reader(bucket, key, constants.RingBufferSize)
}

func (s *Server) GetResults(
	req *api.GetResultsRequest,
	stream api.TestService_GetResultsServer,
) error {
	if len(os.Args[1:]) != 2 {
		log.Fatalf("missing required arguments ./grpcserver bucket key")
	}

	reader := getS3Reader(os.Args[1], os.Args[2])
	buf := make([]byte, constants.S3BufferSize)

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read from buffer: %v\n", err)
		}

		chunk := &api.DataChunk{
			Data: buf,
			Size: int64(n),
		}
		if err := stream.Send(chunk); err != nil {
			log.Fatalf("failed to send chunk: %v\n", err)
		}
	}

	return nil
}

func (s *Server) GetStatus(
	ctx context.Context,
	req *api.GetStatusRequest,
) (*api.GetStatusResponse, error) {
	resp := &api.GetStatusResponse{
		Msg: fmt.Sprintf("Hello %s", req.Msg),
		Ok:  true,
	}
	return resp, nil
}

func main() {
	grpcServer := SetupGrpcServer()
	lis, err := net.Listen("tcp", "0.0.0.0:4040")
	if err != nil {
		log.Fatalf("failed to open TCP socket on port 4040: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("starting gRPC server on port 4040")
	if err := grpcServer.Serve(lis); err != nil {
		if err != grpc.ErrServerStopped {
			log.Fatalf("failed to start gRPC server: %v\n", err)
			os.Exit(1)
		}
	}
}
