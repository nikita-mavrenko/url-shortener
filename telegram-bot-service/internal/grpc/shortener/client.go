package shortener

import (
	"context"
	"fmt"
	shortenerv1 "github.com/nikitamavrenko/telegram-bot-service/proto/shortener"
	"google.golang.org/grpc"
)

type Client struct {
	client shortenerv1.ShortenerClient
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	const op = "grpc.NewShortenerClient"

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := shortenerv1.NewShortenerClient(conn)

	return &Client{
		client: client,
	}, nil
}

func (c *Client) GetShortenLink(ctx context.Context, url string) (string, error) {
	const op = "grpc.ShortURL"

	resp, err := c.client.ShortURL(ctx, &shortenerv1.ShortURLRequest{Url: url})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return resp.ShortenUrl, nil
}
