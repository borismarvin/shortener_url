package handlers

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DatabaseName string

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
)

func CheckDBConn(DatabaseName string) (err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, DatabaseName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return
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
