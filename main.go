package main

import (
	"fmt"
	"log"
	"microservice_job_avito/internal/database"

	"database/sql"

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

	_, err = dotsql.LoadFromFile("start.sql")

	// _, err = dot.Exec(db, "drop-balance-table")

	// if err != nil {
	// 	panic(err)
	// }

	// _, err = dot.Exec(db, "drop-user-table")

	// if err != nil {
	// 	panic(err)
	// }

	// _, err = dot.Exec(db, "create-user-table")

	// if err != nil {
	// 	panic(err)
	// }

	// _, err = dot.Exec(db, "create-balance-table")

	// if err != nil {
	// 	panic(err)
	// }

	// create users if not exist
	// var user uint64

	// row, _ := dot.QueryRow(db, "check-users-exist")
	// err = row.Scan(&user)

	// if err == sql.ErrNoRows {
	// 	// create 3 users

	// 	for i := 0; i < 3; i++ {
	// 		_, err = dot.Exec(db, "create-user")

	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	}

	// } else if err != nil {
	// 	panic(err)
	// }

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

	// os.Exit(0)

	database.PathHandler()
}
