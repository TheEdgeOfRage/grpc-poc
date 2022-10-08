package client

import (
	"log"

	api "grpc-test/gen/proto/go/results/api/v1"
)

func (c *Client) GetStatus() (string, bool) {
	resp, err := c.client.GetStatus(c.ctx, &api.GetStatusRequest{})
	if err != nil {
		log.Fatalf("grpc.GetStatus failed: %v", err)
	}

	return resp.Msg, resp.Ok
}
