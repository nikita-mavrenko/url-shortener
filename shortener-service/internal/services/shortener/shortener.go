package shortener

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nikitamavrenko/shortener-service/internal/domain"
	"github.com/nikitamavrenko/shortener-service/internal/storage"
	"github.com/nikitamavrenko/shortener-service/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
)

type ShortenerService struct {
	log         *zerolog.Logger
	storage     Storage
	baseURL     string
	alphabet    []rune
	alphabetLen uint32
	redis       Redis
}

type Storage interface {
	SaveURL(ctx context.Context, url *domain.URL) error
	GetURL(ctx context.Context, id string) (*domain.URL, error)
}

type Redis interface {
	GetUrl(ctx context.Context, id string) (string, error)
	PutUrl(ctx context.Context, id, url string) error
}

func New(log *zerolog.Logger, storage Storage, redis Redis, baseUrl string) *ShortenerService {
	alphabet := utils.GenerateAlphabet(32)

	return &ShortenerService{
		log:         log,
		storage:     storage,
		baseURL:     baseUrl,
		alphabet:    alphabet,
		alphabetLen: uint32(len(alphabet)),
		redis:       redis,
	}
}

func (s *ShortenerService) Short(ctx context.Context, url string) (string, error) {
	shortenURL := &domain.URL{
		Id:          s.makeIdentifier(),
		OriginalURL: url,
	}

	err := s.storage.SaveURL(ctx, shortenURL)
	if err != nil {
		return "", err
	}

	if err := s.redis.PutUrl(ctx, shortenURL.Id, shortenURL.OriginalURL); err != nil {
		log.Warn().Err(err).Msg("failed to save shorten url to redis")
	}

	return s.makeURL(shortenURL.Id), nil
}

func (s *ShortenerService) Redirect(ctx context.Context, id string) (string, error) {
	url, err := s.redis.GetUrl(ctx, id)
	if err == nil {
		return url, nil
	} else if errors.Is(err, storage.ErrUrlNotFound) {
		s.log.Info().Str("id", id).Str("url", url).Msg("url not found in redis")
	}

	shortenUrl, err := s.storage.GetURL(ctx, id)
	if err != nil {
		return "", err
	}
	return shortenUrl.OriginalURL, nil
}

func (s *ShortenerService) makeIdentifier() string {
	id := uuid.New().ID()

	var (
		indexes = make([]uint32, 0, 10)
		sb      strings.Builder
		num     = id
	)

	for num > 0 {
		indexes = append(indexes, num%s.alphabetLen)
		num /= s.alphabetLen
	}

	for _, index := range indexes {
		sb.WriteString(string(s.alphabet[index]))
	}

	return sb.String()
}

func (s *ShortenerService) makeURL(identifier string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, identifier)
}
