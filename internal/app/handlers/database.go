package handlers

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DSN string

func Ping() error {
	db, err := sql.Open("postgres", DSN) // mysql || postgres
	if err != nil {
		panic(err)
	} else {
		log.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}

// func (d *DBStorage) CheckDBConn(w http.ResponseWriter, r *http.Request) {
// 	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
// 		`localhost`, `video`, `XXXXXXXX`, `video`)
// 	db, err := sql.Open("pgx", ps)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()
// }
