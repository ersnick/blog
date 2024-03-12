package db

import (
	"database/sql"
	"fmt"
)

var Db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1009"
	dbname   = "blogDB"
)

// Connect db
func InitDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = Db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}
