package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func New(
	host string,
	port int,
) (*Client, error) {

	client := redis.NewClient(
		&redis.Options{
			Addr: fmt.Sprintf(
				"%s:%d",
				host,
				port,
			),
		},
	)

	err := client.Ping(
		context.Background(),
	).Err()

	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Client() *redis.Client {
	return c.client
}

func (c *Client) Close() error {

	return c.client.Close()
}
