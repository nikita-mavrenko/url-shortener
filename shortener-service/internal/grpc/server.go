package grpc

import (
	"context"
	shortenerv1 "github.com/nikitamavrenko/shortener-service/proto/shortener"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShortenerServer struct {
	shortenerv1.UnimplementedShortenerServer
	shortener Shortener
}

type Shortener interface {
	Short(ctx context.Context, url string) (string, error)
}

func Register(gRPC *grpc.Server, shortener Shortener) {
	shortenerv1.RegisterShortenerServer(gRPC, &ShortenerServer{shortener: shortener})
}

func (s *ShortenerServer) ShortURL(ctx context.Context, in *shortenerv1.ShortURLRequest) (*shortenerv1.ShortURLResponse, error) {
	if in.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "url is required")
	}

	shortenURL, err := s.shortener.Short(ctx, in.Url)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &shortenerv1.ShortURLResponse{
		ShortenUrl: shortenURL,
	}, nil
}
