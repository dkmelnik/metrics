package collect

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dkmelnik/metrics/proto/metrics"
)

type GRPCClient struct {
	metrics.MetricsClient
	conn *grpc.ClientConn
}

func NewGRPCClient(addr string) (*GRPCClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GRPCClient{metrics.NewMetricsClient(conn), conn}, nil
}

type Payload struct {
	ID    string
	MType string
	Delta int64
	Value float64
}

func (c *GRPCClient) Send(ctx context.Context, p Payload) error {
	_, err := c.Create(ctx, &metrics.CreateRequest{
		Id:    p.ID,
		Mtype: p.MType,
		Delta: p.Delta,
		Value: p.Value,
	})
	if err != nil {
		return err
	}
	return nil
}
