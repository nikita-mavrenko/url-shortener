package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/nikitamavrenko/shortener-service/internal/domain"
	"github.com/nikitamavrenko/shortener-service/internal/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"runtime"
	"time"
)

type Storage struct {
	pool *pgxpool.Pool
	log  *zerolog.Logger
}

func New(ctx context.Context, dbUrl string, log *zerolog.Logger) (*Storage, error) {
	const op = "storage.postgres.New"

	pool, err := connectWithRetry(ctx, dbUrl, 10)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{pool: pool, log: log}, nil
}

func (s *Storage) SaveURL(ctx context.Context, url *domain.URL) error {
	const op = "storage.postgres.SaveURL"

	query := "INSERT INTO urls (id, original_url) VALUES ($1, $2)"
	_, err := s.pool.Exec(ctx, query, url.Id, url.OriginalURL)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("%s: %w", op, storage.ErrUrlAlreadyExists)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetURL(ctx context.Context, id string) (*domain.URL, error) {
	const op = "storage.postgres.GetURL"

	query := "SELECT id, original_url FROM urls WHERE id = $1"
	url := domain.URL{}
	err := s.pool.QueryRow(ctx, query, id).Scan(&url.Id, &url.OriginalURL)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.NoData {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrUrlNotFound)
			}
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &url, nil
}

func (s *Storage) Migrate(ctx context.Context) error {
	const op = "storage.postgres.RunMigrations"
	s.log.Info().Msg("starting migration")

	sqlDb := stdlib.OpenDBFromPool(s.pool)
	defer sqlDb.Close()

	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	migrationsDir := filepath.Join(baseDir, "../../../migrations")

	driver, err := pgx.WithInstance(sqlDb, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://"+migrationsDir, "postgres", driver)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer migrator.Close()

	version, dirty, err := migrator.Version()
	s.log.Info().Int("version", int(version))
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("%s: %w", op, err)
	}

	if dirty {
		s.log.Warn().Int("version", int(version)).Msg("fixing dirty database state")
		err = migrator.Force(int(version))
		if err != nil {
			return fmt.Errorf("%s: force version failed: %w", op, err)
		}
	}

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info().Msg("migrated")
	return nil
}

func connectWithRetry(ctx context.Context, dbUrl string, maxAttempts int) (*pgxpool.Pool, error) {
	var err error
	var pool *pgxpool.Pool
	log.Info().Msg("starting db connection")
	for attempt := 0; attempt < maxAttempts; attempt++ {
		pool, err = pgxpool.New(ctx, dbUrl)
		if err != nil {
			log.Warn().Err(err).Msgf("attempt %d failed", attempt+1)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if err = pool.Ping(ctx); err != nil {
			pool.Close()
			log.Warn().Err(err).Msgf("ping attempt %d failed", attempt+1)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		log.Info().Msg("db connected")
		return pool, nil
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxAttempts, err)
}
