package client

import (
	"context"
	"log"
	"math"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "grpc-test/gen/proto/go/results/api/v1"
)

type Client struct {
	client api.TestServiceClient
	conn   *grpc.ClientConn
	ctx    context.Context
}

func NewClient(ctx context.Context) *Client {
	conn, err := grpc.Dial(
		"127.0.0.1:4040",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32),
			grpc.MaxCallSendMsgSize(math.MaxInt32),
		),
	)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return &Client{
		client: api.NewTestServiceClient(conn),
		conn:   conn,
		ctx:    ctx,
	}
}

func (c *Client) Close() {
	c.conn.Close()
}
