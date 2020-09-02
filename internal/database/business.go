package database

import (
	"log"
	"microservice_job_avito/internal/types"
	"net/http"

	"database/sql"
	"encoding/json"
	"fmt"

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

func _getBalance(userID string) ([]byte, int) {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		js, _ := json.Marshal(NoDatabaseConnection500rm)
		return js, http.StatusInternalServerError
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	// check if user exists
	userRow, err := dot.QueryRow(db, "check-user-exists", userID)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	var userIDCheck string
	err = userRow.Scan(&userIDCheck)

	if err == sql.ErrNoRows {
		js, _ := json.Marshal(UserDoesNotExist400rm)
		return js, http.StatusBadRequest
	} else if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	// if user exists
	b, _ := dot.QueryRow(db, "get-user-balance", userID)

	var balance float64
	err = b.Scan(&balance)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	returnBalance := types.Balance{Bal: balance}

	js, _ := json.Marshal(returnBalance)
	return js, http.StatusOK
}

func _increase(userID string, money float64) ([]byte, int) {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		js, _ := json.Marshal(NoDatabaseConnection500rm)
		return js, http.StatusInternalServerError
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	// check if user exists
	userRow, err := dot.QueryRow(db, "check-user-exists", userID)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	var userIDCheck string
	err = userRow.Scan(&userIDCheck)

	if err == sql.ErrNoRows {

		_, err = dot.Exec(db, "create-user", userID)

		if err != nil {
			js, _ := json.Marshal(InternalServerError500rm)
			return js, http.StatusInternalServerError
		}

		_, err = dot.Exec(db, "create-balance", userID)

		if err != nil {
			js, _ := json.Marshal(InternalServerError500rm)
			return js, http.StatusInternalServerError
		}

	} else if err != nil {
		log.Fatal(err)
	}

	// increase balance

	_, err = dot.Exec(db, "remittance-to", money, userID)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	b, err := dot.QueryRow(db, "get-user-balance", userID)

	var balance float64
	err = b.Scan(&balance)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	returnBalance := types.Balance{Bal: balance}

	js, _ := json.Marshal(returnBalance)
	return js, http.StatusOK
}

func _decrease(userID string, money float64) {
	//
}

func _remittance(userIDFrom string, userIDTo string, money float64) {
	//
}
