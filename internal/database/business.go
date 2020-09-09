package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"microservice_job_avito/internal/types"
	"net/http"

	"github.com/gchaincl/dotsql"
	"github.com/lib/pq"
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

func _getBalance(userID string, args ...interface{}) ([]byte, int, interface{}) {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		js, _ := json.Marshal(NoDatabaseConnection500rm)
		return js, http.StatusInternalServerError, NoDatabaseConnection500rm["error"]["status_message"]
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	// Проверить, существует ли пользователь с данным id
	userRow, err := dot.QueryRow(db, "check-user-exists", userID)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	var userIDCheck string
	err = userRow.Scan(&userIDCheck)

	if err == sql.ErrNoRows {
		js, _ := json.Marshal(UserDoesNotExist400rm)
		return js, http.StatusBadRequest, UserDoesNotExist400rm["error"]["status_message"]
	} else if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	// Если пользователь существует
	b, _ := dot.QueryRow(db, "get-user-balance", userID)

	var balance float64
	err = b.Scan(&balance)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	if flag := args[0].(bool); flag == false {
		returnBalance := types.Balance{Bal: balance}

		js, _ := json.Marshal(returnBalance)
		return js, http.StatusOK, "OK"
	}

	base := args[1].(string)
	currencyReq := fmt.Sprintf("https://api.exchangeratesapi.io/latest?base=%s", base)

	resp, err := http.Get(currencyReq)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	var c types.Currency
	err = json.Unmarshal(body, &c)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	if len(c.Rates) == 0 {
		js, _ := json.Marshal(WrongBaseFormat400rm)
		return js, http.StatusBadRequest, WrongBaseFormat400rm["error"]["status_message"]
	}

	balance /= c.Rates["RUB"]
	balance = math.Round(balance*100) / 100

	returnBalance := types.Balance{Bal: balance}

	js, _ := json.Marshal(returnBalance)
	return js, http.StatusOK, "OK"
}

func _increase(userID string, money float64) ([]byte, int, interface{}) {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		js, _ := json.Marshal(NoDatabaseConnection500rm)
		return js, http.StatusInternalServerError, NoDatabaseConnection500rm["error"]["status_message"]
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	// Проверить, существует ли пользователь с данным id
	userRow, err := dot.QueryRow(db, "check-user-exists", userID)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	var userIDCheck string
	err = userRow.Scan(&userIDCheck)

	if err == sql.ErrNoRows {

		if _, err = dot.Exec(db, "create-user", userID); err != nil {
			js, _ := json.Marshal(InternalServerError500rm)
			return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
		}

		if _, err = dot.Exec(db, "create-balance", userID); err != nil {
			js, _ := json.Marshal(InternalServerError500rm)
			return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
		}

	} else if err != nil {
		log.Fatal(err)
	}

	// Увелечить баланс
	_, err = dot.Exec(db, "remittance-to", money, userID)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	b, err := dot.QueryRow(db, "get-user-balance", userID)

	var balance float64
	err = b.Scan(&balance)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	returnBalance := types.Balance{Bal: balance}

	js, _ := json.Marshal(returnBalance)
	return js, http.StatusOK, "OK"
}

func _decrease(userID string, money float64) ([]byte, int, interface{}) {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		js, _ := json.Marshal(NoDatabaseConnection500rm)
		return js, http.StatusInternalServerError, NoDatabaseConnection500rm["error"]["status_message"]
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	// Проверить, существует ли пользователь с данным id
	userRow, err := dot.QueryRow(db, "check-user-exists", userID)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	var userIDCheck string
	err = userRow.Scan(&userIDCheck)

	if err == sql.ErrNoRows {
		js, _ := json.Marshal(UserDoesNotExist400rm)
		return js, http.StatusBadRequest, UserDoesNotExist400rm["error"]["status_message"]
	} else if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	// Уменьшить баланс
	b, err := dot.QueryRow(db, "get-user-balance", userID)

	var balance float64
	err = b.Scan(&balance)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	} else if balance < money { // Пользовтелеь не имеет достаточное количество средств для списания
		js, _ := json.Marshal(DecreaseMore400rm)
		return js, http.StatusBadRequest, DecreaseMore400rm["error"]["status_message"]
	}

	if _, err = dot.Exec(db, "remittance-from", money, userID); err != nil {

		if err, _ := err.(*pq.Error); err.Code.Name() == "check_violation" {
			js, _ := json.Marshal(DecreaseMore400rm)
			return js, http.StatusBadRequest, DecreaseMore400rm["error"]["status_message"]
		}

		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	b, err = dot.QueryRow(db, "get-user-balance", userID)
	err = b.Scan(&balance)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	returnBalance := types.Balance{Bal: balance}

	js, _ := json.Marshal(returnBalance)
	return js, http.StatusOK, "OK"
}

func _remittance(userIDFrom string, userIDTo string, money float64) ([]byte, int, interface{}) {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		js, _ := json.Marshal(NoDatabaseConnection500rm)
		return js, http.StatusInternalServerError, NoDatabaseConnection500rm["error"]["status_message"]
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	// Проверить, существует ли пользователь, с баланса которого нужно списать средства
	userRow, err := dot.QueryRow(db, "check-user-exists", userIDFrom)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	var userIDCheck string
	err = userRow.Scan(&userIDCheck)

	if err == sql.ErrNoRows {
		js, _ := json.Marshal(UserDoesNotExist400rm)
		return js, http.StatusBadRequest, UserDoesNotExist400rm["error"]["status_message"]
	} else if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	// Проверить, существует ли пользователь, на баланс которому нужно перевести средства
	userRow, err = dot.QueryRow(db, "check-user-exists", userIDTo)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	err = userRow.Scan(&userIDCheck)

	if err == sql.ErrNoRows {
		js, _ := json.Marshal(UserDoesNotExist400rm)
		return js, http.StatusBadRequest, UserDoesNotExist400rm["error"]["status_message"]
	} else if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	// Проверить, что пользователь, с баланса которого нужно списать средства, имеет достаточно средств для списания
	// b, err := dot.QueryRow(db, "get-user-balance", userIDFrom)

	// var balance float64
	// err = b.Scan(&balance)

	// if err != nil {
	// 	js, _ := json.Marshal(InternalServerError500rm)
	// 	return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	// } else if balance < money {
	// 	js, _ := json.Marshal(DecreaseMore400rm)
	// 	return js, http.StatusBadRequest, DecreaseMore400rm["error"]["status_message"]
	// }

	if _, err = dot.Exec(db, "remittance-from", money, userIDFrom); err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	if _, err = dot.Exec(db, "remittance-to", money, userIDTo); err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	var balanceFrom, balanceTo float64

	bFrom, err := dot.QueryRow(db, "get-user-balance", userIDFrom)
	err = bFrom.Scan(&balanceFrom)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	bTo, err := dot.QueryRow(db, "get-user-balance", userIDTo)
	err = bTo.Scan(&balanceTo)

	if err != nil {
		js, _ := json.Marshal(InternalServerError500rm)
		return js, http.StatusInternalServerError, InternalServerError500rm["error"]["status_message"]
	}

	returnBalance := types.Remittance{BalanceFrom: balanceFrom, BalanceTo: balanceTo}

	js, _ := json.Marshal(returnBalance)
	return js, http.StatusOK, "OK"
}
