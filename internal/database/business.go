package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"microservice_job_avito/internal/types"
	"net/http"

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

	// Проверить, существует ли пользователь с данным id
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

	// Если пользователь существует
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

	// Проверить, существует ли пользователь с данным id
	userRow, err := dot.QueryRow(db, "check-user-exists", userID)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	var userIDCheck string
	err = userRow.Scan(&userIDCheck)

	if err == sql.ErrNoRows {

		if _, err = dot.Exec(db, "create-user", userID); err != nil {
			js, _ := json.Marshal(InternalServerError500rm)
			return js, http.StatusInternalServerError
		}

		if _, err = dot.Exec(db, "create-balance", userID); err != nil {
			js, _ := json.Marshal(InternalServerError500rm)
			return js, http.StatusInternalServerError
		}

	} else if err != nil {
		log.Fatal(err)
	}

	// Увелечить баланс
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

func _decrease(userID string, money float64) ([]byte, int) {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		js, _ := json.Marshal(NoDatabaseConnection500rm)
		return js, http.StatusInternalServerError
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	// Проверить, существует ли пользователь с данным id
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

	// Уменьшить баланс
	b, err := dot.QueryRow(db, "get-user-balance", userID)

	var balance float64
	err = b.Scan(&balance)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	} else if balance < money { // Пользовтелеь не имеет достаточное количество средств для списания
		js, _ := json.Marshal(DecreaseMore400rm)
		return js, http.StatusBadRequest
	}

	if _, err = dot.Exec(db, "remittance-from", money, userID); err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	b, err = dot.QueryRow(db, "get-user-balance", userID)
	err = b.Scan(&balance)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	returnBalance := types.Balance{Bal: balance}

	js, _ := json.Marshal(returnBalance)
	return js, http.StatusOK
}

func _remittance(userIDFrom string, userIDTo string, money float64) ([]byte, int) {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		js, _ := json.Marshal(NoDatabaseConnection500rm)
		return js, http.StatusInternalServerError
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	// Проверить, существует ли пользователь, с баланса которого нужно списать средства
	userRow, err := dot.QueryRow(db, "check-user-exists", userIDFrom)

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

	// Проверить, существует ли пользователь, на баланс которому нужно перевести средства
	userRow, err = dot.QueryRow(db, "check-user-exists", userIDTo)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	err = userRow.Scan(&userIDCheck)

	if err == sql.ErrNoRows {
		js, _ := json.Marshal(UserDoesNotExist400rm)
		return js, http.StatusBadRequest
	} else if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	// Проверить, что пользователь, с баланса которого нужно списать средства, имеет достаточно средств для списания
	b, err := dot.QueryRow(db, "get-user-balance", userIDFrom)

	var balance float64
	err = b.Scan(&balance)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	} else if balance < money {
		js, _ := json.Marshal(DecreaseMore400rm)
		return js, http.StatusBadRequest
	}

	if _, err = dot.Exec(db, "remittance-from", money, userIDFrom); err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	if _, err = dot.Exec(db, "remittance-to", money, userIDTo); err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	var balanceFrom, balanceTo float64

	bFrom, err := dot.QueryRow(db, "get-user-balance", userIDFrom)
	err = bFrom.Scan(&balanceFrom)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	bTo, err := dot.QueryRow(db, "get-user-balance", userIDTo)
	err = bTo.Scan(&balanceTo)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError
	}

	returnBalance := types.Remittance{BalanceFrom: balanceFrom, BalanceTo: balanceTo}

	js, _ := json.Marshal(returnBalance)
	return js, http.StatusOK
}
