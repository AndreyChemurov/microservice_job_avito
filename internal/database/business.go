package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"microservice_job_avito/internal/currency"
	"microservice_job_avito/internal/types"

	"github.com/gchaincl/dotsql"
	"github.com/lib/pq"
)

const (
	host     = "db"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

var postgresInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

func getBalance(userID string, flag bool, base string) (*types.Balance, error) {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		response := types.Balance{Status: 500}

		return &response, errors.New("Service cannot connect to database")
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	// Проверить, существует ли пользователь с данным id
	userRow, err := dot.QueryRow(db, "check-user-exists", userID)

	if err != nil {
		response := types.Balance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	var userIDCheck string
	err = userRow.Scan(&userIDCheck)

	if err == sql.ErrNoRows {
		response := types.Balance{Status: 400}

		return &response, errors.New("User does not exist")

	} else if err != nil {
		response := types.Balance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	// Если пользователь существует
	response, _ := dot.QueryRow(db, "get-user-balance", userID)

	var balance float64
	err = response.Scan(&balance)

	if err != nil {
		response := types.Balance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	if flag == false {
		returnBalance := types.Balance{Balance: balance, Status: 200}

		return &returnBalance, nil
	}

	balance, err = currency.Currency(balance, base)

	if err != nil {
		response := types.Balance{Status: 500}

		return &response, err
	}

	returnBalance := types.Balance{Balance: balance, Status: 200}

	return &returnBalance, nil
}

func increase(userID string, money float64) (*types.Balance, error) {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		response := types.Balance{Status: 500}

		return &response, errors.New("Service cannot connect to database")
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	// Проверить, существует ли пользователь с данным id
	userRow, err := dot.QueryRow(db, "check-user-exists", userID)

	if err != nil {
		response := types.Balance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	var userIDCheck string
	err = userRow.Scan(&userIDCheck)

	if err == sql.ErrNoRows {

		if _, err = dot.Exec(db, "create-user", userID); err != nil {
			response := types.Balance{Status: 500}

			return &response, errors.New("Internal server error")
		}

		if _, err = dot.Exec(db, "create-balance", userID); err != nil {
			response := types.Balance{Status: 500}

			return &response, errors.New("Internal server error")
		}

	} else if err != nil {
		log.Fatal(err)
	}

	// Увелечить баланс
	_, err = dot.Exec(db, "remittance-to", money, userID)

	if err != nil {
		response := types.Balance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	response, err := dot.QueryRow(db, "get-user-balance", userID)

	var balance float64
	err = response.Scan(&balance)

	if err != nil {
		response := types.Balance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	returnBalance := types.Balance{Balance: balance, Status: 200}
	return &returnBalance, nil
}

func decrease(userID string, money float64) (*types.Balance, error) {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		response := types.Balance{Status: 500}

		return &response, errors.New("Service cannot connect to database")
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	// Проверить, существует ли пользователь с данным id
	userRow, err := dot.QueryRow(db, "check-user-exists", userID)

	if err != nil {
		response := types.Balance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	var userIDCheck string
	err = userRow.Scan(&userIDCheck)

	if err == sql.ErrNoRows {
		response := types.Balance{Status: 500}

		return &response, errors.New("User does not exist")

	} else if err != nil {
		response := types.Balance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	// Уменьшить баланс
	response, err := dot.QueryRow(db, "get-user-balance", userID)

	var balance float64
	err = response.Scan(&balance)

	if err != nil {
		response := types.Balance{Status: 500}

		return &response, errors.New("Internal server error")

	} else if balance < money { // Пользовтелеь не имеет достаточное количество средств для списания
		response := types.Balance{Status: 500}

		return &response, errors.New("Operation not available: Debit exceeds the balance")
	}

	if _, err = dot.Exec(db, "remittance-from", money, userID); err != nil {

		if err, _ := err.(*pq.Error); err.Code.Name() == "check_violation" {
			response := types.Balance{Status: 500}

			return &response, errors.New("Operation not available: Debit exceeds the balance")
		}

		response := types.Balance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	response, err = dot.QueryRow(db, "get-user-balance", userID)
	err = response.Scan(&balance)

	if err != nil {
		response := types.Balance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	returnBalance := types.Balance{Balance: balance, Status: 200}
	return &returnBalance, nil
}

func remittance(userIDFrom string, userIDTo string, money float64) (*types.Remittance, error) {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		response := types.Remittance{Status: 500}

		return &response, errors.New("Service cannot connect to database")
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	// Проверить, существует ли пользователь, с баланса которого нужно списать средства
	userRow, err := dot.QueryRow(db, "check-user-exists", userIDFrom)

	if err != nil {
		response := types.Remittance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	var userIDCheck string
	err = userRow.Scan(&userIDCheck)

	if err == sql.ErrNoRows {
		response := types.Remittance{Status: 500}

		return &response, errors.New("User does not exist")

	} else if err != nil {
		response := types.Remittance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	// Проверить, существует ли пользователь, на баланс которому нужно перевести средства
	userRow, err = dot.QueryRow(db, "check-user-exists", userIDTo)

	if err != nil {
		response := types.Remittance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	err = userRow.Scan(&userIDCheck)

	if err == sql.ErrNoRows {
		response := types.Remittance{Status: 500}

		return &response, errors.New("User does not exist")

	} else if err != nil {
		response := types.Remittance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	if _, err = dot.Exec(db, "remittance-from", money, userIDFrom); err != nil {
		response := types.Remittance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	if _, err = dot.Exec(db, "remittance-to", money, userIDTo); err != nil {
		response := types.Remittance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	var balanceFrom, balanceTo float64

	bFrom, err := dot.QueryRow(db, "get-user-balance", userIDFrom)
	err = bFrom.Scan(&balanceFrom)

	if err != nil {
		response := types.Remittance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	bTo, err := dot.QueryRow(db, "get-user-balance", userIDTo)
	err = bTo.Scan(&balanceTo)

	if err != nil {
		response := types.Remittance{Status: 500}

		return &response, errors.New("Internal server error")
	}

	returnBalance := types.Remittance{BalanceFrom: balanceFrom, BalanceTo: balanceTo, Status: 200}

	return &returnBalance, nil
}
