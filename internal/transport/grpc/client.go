package controlplanegrpc

import (
	"context"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	addr string
	mu   sync.Mutex
	conn *grpc.ClientConn
}

// New creates a reusable gRPC client wrapper for outbound adminui -> controlplane calls.
func New(addr string) *Client {
	return &Client{addr: strings.TrimSpace(addr)}
}

// Invoke reuses one gRPC connection and wraps each outbound call with a short timeout.
func (c *Client) Invoke(ctx context.Context, fn func(AdminCatalogServiceClient) error) error {
	callCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	conn, err := c.connection(callCtx)
	if err != nil {
		return err
	}
	return fn(NewAdminCatalogServiceClient(conn))
}

// connection lazily dials controlplane once and then reuses the same client connection.
func (c *Client) connection(ctx context.Context) (*grpc.ClientConn, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c.conn, nil
	}
	conn, err := grpc.DialContext(
		ctx,
		c.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	return conn, nil
}
