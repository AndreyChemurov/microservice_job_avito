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

// 1. check if user exists
// 2. get balance
// 3. if not exists, return error
func _getBalance(userID uint64) error {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	balance, err := dot.QueryRow(db, "get-user-balance", userID)

	if err == sql.ErrNoRows || err != nil {
		panic(err)
	}

	log.Println(balance) // return json
	return nil
}

// 1. check if user exists
// 2. increase balance
// 3. if not exists, return error
func _increase(userID uint64, money float64) {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	_, err = dot.QueryRow(db, "get-user-balance", userID)

	if err == sql.ErrNoRows {
		_, err = dot.Exec(db, "create-balance", userID)

		if err != nil {
			panic(err)
		}

	} else if err != nil {
		panic(err)
	}

	_, err = dot.Exec(db, "remittance-to", money, userID)

	if err != nil {
		panic(err)
	}
}

// 1. check if user exists
// 2. decrease balance
// 3. if not exists, return error
func _decrease(userID uint64, money float64) {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	_, err = dot.QueryRow(db, "get-user-balance", userID)

	if err == sql.ErrNoRows || err != nil {
		panic(err)
	}

	_, err = dot.Exec(db, "remittance-from", money, userID)

	if err != nil {
		panic(err)
	}
}

// 1. check if both users exist
// 2. decrease balance from user 1
// 3. increase balance for user 2
func _remittance(userIDFrom uint64, userIDTo uint64, money float64) {
	//
}
