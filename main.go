package main

import (
	"database/sql"
	"fmt"
	"log"
	"microservice_job_avito/internal/database"

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

func prepareDatabase() error {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	if err != nil {
		panic(err)
	}

	// _, err = dot.Exec(db, "drop-balance-table")

	// if err != nil {
	// 	panic(err)
	// }

	// _, err = dot.Exec(db, "drop-user-table")

	// if err != nil {
	// 	panic(err)
	// }

	_, err = dot.Exec(db, "create-user-table")

	if err != nil {
		panic(err)
	}

	_, err = dot.Exec(db, "create-balance-table")

	if err != nil {
		panic(err)
	}

	return nil
}

func main() {
	err := prepareDatabase()

	if err != nil {
		log.Fatal(err)
	}

	database.PathHandler()
}
