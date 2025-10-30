package shortener

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nikitamavrenko/shortener-service/internal/domain"
	"github.com/nikitamavrenko/shortener-service/internal/utils"
	"github.com/rs/zerolog"
	"strings"
)

type ShortenerService struct {
	log         *zerolog.Logger
	storage     Storage
	baseURL     string
	alphabet    []rune
	alphabetLen uint32
}

type Storage interface {
	SaveURL(ctx context.Context, url *domain.URL) error
	GetURL(ctx context.Context, id string) (*domain.URL, error)
}

func New(log *zerolog.Logger, storage Storage, baseUrl string) *ShortenerService {
	alphabet := utils.GenerateAlphabet(32)

	return &ShortenerService{
		log:         log,
		storage:     storage,
		baseURL:     baseUrl,
		alphabet:    alphabet,
		alphabetLen: uint32(len(alphabet)),
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

	return s.makeURL(shortenURL.Id), nil
}

func (s *ShortenerService) Redirect(ctx context.Context, id string) (string, error) {
	url, err := s.storage.GetURL(ctx, id)
	if err != nil {
		return "", err
	}
	return url.OriginalURL, nil
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
