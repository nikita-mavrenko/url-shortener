package redis

import (
	"context"
	"errors"
	"github.com/nikitamavrenko/shortener-service/internal/config"
	"github.com/nikitamavrenko/shortener-service/internal/storage"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"time"
)

type Redis struct {
	client *redis.Client
	log    *zerolog.Logger
}

func New(ctx context.Context, log *zerolog.Logger, cfg *config.Config) (*Redis, error) {
	log.Info().Msg("init redis")
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	log.Info().Msg("init redis success")
	return &Redis{
		client: client,
		log:    log,
	}, nil
}

func (r *Redis) GetUrl(ctx context.Context, id string) (string, error) {
	r.log.Info().Str("id", id).Msg("get url from redis")
	url, err := r.client.Get(ctx, id).Result()
	if err != nil {
		r.log.Error().Err(err).Str("id", id).Msg("failed to get url from redis")
		if errors.Is(err, redis.Nil) {
			return "", storage.ErrUrlNotFound
		}
		return "", err
	}
	return url, nil
}

func (r *Redis) PutUrl(ctx context.Context, id, url string) error {
	r.log.Info().Str("id", id).Str("url", url).Msg("put url into redis")

	if _, err := r.client.Set(ctx, id, url, time.Hour*24).Result(); err != nil {
		r.log.Error().Err(err).Str("id", id).Msg("failed to put url into redis")
		return err
	}
	return nil
}
