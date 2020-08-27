package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gchaincl/dotsql"
	_ "github.com/lib/pq" //
)

const (
	host     = "db"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

var postgresInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

func _getBalance(userID uint64) error {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	balance, err := dot.Exec(db, "get-user-balance")

	if err == sql.ErrNoRows || err != nil {
		panic(err)
	}

	log.Println(balance) // return json
	return nil
}

func _increase(usedID uint64, money float64) {
	//
}

func _decrease(usedID uint64, money float64) {
	//
}

func _remittance(userIDFrom uint64, userIDTo uint64, money float64) {
	//
}
