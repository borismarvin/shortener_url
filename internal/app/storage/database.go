package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	//_ "github.com/go-sql-driver/mysql"
	"github.com/borismarvin/shortener_url.git/internal/app/types"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBRepository struct {
	DB  *sqlx.DB
	cfg *types.Config
}

func NewDBRepository(cfg *types.Config) *DBRepository {
	repo := &DBRepository{
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

func (r *DBRepository) Save(url *types.URL) (err error) {
	if r.DB == nil {
		fmt.Println("ошибка подлкючения к бд")
	}
	rows, err := r.DB.QueryContext(context.Background(), "SELECT * FROM urls where 'hash' = $1", url.Hash)

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	if rows.Err() != nil {
		return err
	}

	if rows.Next() {
		fmt.Println("конфликт url")
	}

	_, err = r.DB.NamedExec(`INSERT INTO urls (hash, uuid, url, short_url)
		VALUES (:hash, :uuid, :url, :short_url)`, url)

	return err
}

func (r *DBRepository) FindByHash(hash string) (exist bool, url *types.URL, err error) {
	if r.DB == nil {
		exist = false
		url = nil
		err = errors.New("нет подключения к бд")
		return
	}

	rows, err := r.DB.QueryContext(context.Background(), "SELECT u.hash, u.uuid, u.url, u.short_url FROM urls u WHERE u.hash = $1 limit $2", hash, 1)

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	if rows.Err() != nil {
		exist = false
		return
	}

	url = &types.URL{}
	for rows.Next() {
		exist = true
		rows.Scan(&url.Hash, &url.UUID, &url.URL, &url.ShortURL)
	}

	return
}

func (r *DBRepository) Ping() (err error) {
	if r.DB == nil {
		return errors.New("нет подключения к бд")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.DB.PingContext(ctx)
}

func (r *DBRepository) migrate() {
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
