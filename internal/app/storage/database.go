package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/borismarvin/shortener_url.git/internal/app/types"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DSN string

type DatabaseRepository struct {
	DB  *sqlx.DB
	cfg *types.Config
}

func NewDatabaseRepository(cfg *types.Config) *DatabaseRepository {
	repo := &DatabaseRepository{
		cfg: cfg,
		DB:  nil,
	}

	if cfg.DatabaseDsn != "" {
		db, err := sqlx.Open("postgres", cfg.DatabaseDsn) // mysql || postgres
		if err == nil {
			repo.DB = db
			repo.migrate()
		} else {
			log.Println(err)
		}
	}

	return repo
}

func (r *DatabaseRepository) Ping() error {
	if r.DB == nil {
		return errors.New("нет подключения к бд")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return r.DB.PingContext(ctx)
}
func (r *DatabaseRepository) Save(url *types.URL) (err error) {
	if r.DB == nil {
		return fmt.Errorf("no connection")
	}

	var existingURL types.URL
	err = r.DB.GetContext(context.Background(), &existingURL, "SELECT * FROM urls WHERE hash = $1", url.Hash)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if existingURL.Hash != "" {
		return fmt.Errorf("URL conflict")
	}

	_, err = r.DB.NamedExec(`INSERT INTO urls (hash, uuid, url, short_url)
        VALUES (:hash, :uuid, :url, :short_url)`, url)
	return err
}
func (r *DatabaseRepository) FindByHash(hash string) (exist bool, url *types.URL, err error) {
	if r.DB == nil {
		err = errors.New("no connection to the database")
		return
	}

	var urls []*types.URL
	err = r.DB.SelectContext(context.Background(), &urls, "SELECT hash, uuid, url, short_url FROM urls WHERE hash = $1 LIMIT 1", hash)
	if err != nil {
		return
	}

	if len(urls) == 0 {
		exist = false
		return
	}

	exist = true
	url = urls[0]
	return
}
func (r *DatabaseRepository) SaveBatch(url []*types.URL) (err error) {
	if r.DB == nil {
		err = errors.New("нет подключения к бд")
		return
	}
	for u := range url {
		_, err = r.DB.NamedExec(`INSERT INTO urls (hash, uuid, url, short_url)
        VALUES (:hash, :uuid, :url, :short_url)`, u)
	}

	return err
}
func (r *DatabaseRepository) migrate() {
	_, err := r.DB.Exec(`CREATE TABLE IF NOT EXISTS urls
		(
			hash      varchar(256) not null,
			uuid      varchar(256) not null,
			url       text         not null,
			short_url varchar(256) not null,
			constraint uk
				unique (hash, uuid)
		)`,
	)

	log.Println(err)
}
