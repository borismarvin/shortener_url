package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/borismarvin/shortener_url.git/internal/app/entity"
	"github.com/borismarvin/shortener_url.git/internal/app/storage/api/model"
	"github.com/borismarvin/shortener_url.git/internal/app/storage/postgres/migration"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	model.Storage

	db *sql.DB
}

func NewPostgresStorage(dbStorageConnect string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dbStorageConnect)
	if err != nil {
		return nil, fmt.Errorf("error while postgresql connect: %w", err)
	}

	err = migration.InitDBTables(db)
	if err != nil {
		return nil, fmt.Errorf("error while postgresql table initialization, %w", err)
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) Close() entity.Response {
	err := s.db.Close()
	if err != nil {
		outErr := fmt.Errorf("couldn'r closed postgres db: %w", err)
		return entity.ErrorResponse(outErr)
	}

	return entity.OKResponse()
}

func (s *PostgresStorage) PingServer(ctx context.Context) entity.Response {
	err := s.db.PingContext(ctx)
	if err != nil {
		outErr := fmt.Errorf("couldn't ping postgres server: %w", err)
		return entity.ErrorResponse(outErr)
	}

	return entity.OKResponse()
}

func (s *PostgresStorage) AddURL(ctx context.Context, key, value entity.URL) entity.Response {
	query := `INSERT INTO url(short_url, url) VALUES(@shortUrl, @url)`
	args := pgx.NamedArgs{
		"shortUrl": key.String(),
		"url":      value.String(),
	}

	_, err := s.db.ExecContext(ctx, query, args)
	if err != nil {
		return entity.ErrorResponse(fmt.Errorf("unable to insert row to postgres: %w", err))
	}

	return entity.OKResponse()
}

func (s *PostgresStorage) GetURL(ctx context.Context, key entity.URL) entity.URLResponse {
	query := `SELECT url FROM url WHERE short_url=@shortUrl`
	args := pgx.NamedArgs{
		"shortUrl": key.String(),
	}

	var dbURL string
	row := s.db.QueryRowContext(ctx, query, args)
	if row == nil {
		return entity.ErrorURLResponse(fmt.Errorf("error while postgres request execution"))
	}

	err := row.Scan(&dbURL)
	if err != nil {
		return entity.ErrorURLResponse(fmt.Errorf("error while processing response row in postgres: %w", err))
	}

	url, err := entity.NewURL(dbURL)
	if err != nil {
		return entity.ErrorURLResponse(fmt.Errorf("error while creating url in postgres: %w", err))
	}

	return entity.OKURLResponse(*url)
}
