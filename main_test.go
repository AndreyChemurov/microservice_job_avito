package main

import (
	"database/sql"
	"fmt"
	"log"
	"microservice_job_avito/internal/database"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	constURL string = "http://localhost:8000"
)

var (
	url string = ""
)

type function func(http.ResponseWriter, *http.Request)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	err := prepareDatabase()

	if err != nil {
		log.Fatalln(err)
	}

	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO user_job VALUES ($1)", "test-user-1")

	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec("INSERT INTO user_job VALUES ($1)", "test-user-2")

	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec("INSERT INTO balance_job VALUES (DEFAULT, $1, 0)", "test-user-1")

	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec("INSERT INTO balance_job VALUES (DEFAULT, $1, 0)", "test-user-2")

	if err != nil {
		log.Fatalln(err)
	}
}

func teardown() {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM balance_job WHERE user_id = $1", "test-user-1")

	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec("DELETE FROM balance_job WHERE user_id = $1", "test-user-2")

	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec("DELETE FROM user_job WHERE user_id = $1", "test-user-1")

	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec("DELETE FROM user_job WHERE user_id = $1", "test-user-2")

	if err != nil {
		log.Fatalln(err)
	}
}

func testingTemplate(route string, callingfn function, params map[string]string) error {
	url = (constURL + route + "?")

	for k, v := range params {
		url += (k + "=" + v + "&")
	}
	url = url[:len(url)-1]

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	rec := httptest.NewRecorder()

	callingfn(rec, request)

	if want, got := 200, rec.Code; want != got {
		return fmt.Errorf("want %d, got %d", want, got)
	}

	return nil
}

func TestGetBalance(t *testing.T) {
	data := map[string]string{
		"id": "test-user-1",
	}

	err := testingTemplate("/balance", database.GetBalance, data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGetBalanceCurrency(t *testing.T) {
	data := map[string]string{
		"id":       "test-user-1",
		"currency": "USD",
	}

	err := testingTemplate("/balance", database.GetBalance, data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestIncreaseUser1(t *testing.T) {
	data := map[string]string{
		"id":    "test-user-1",
		"money": "1000.10",
	}

	err := testingTemplate("/increase", database.IncreaseAndDecrease, data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestIncreaseUser2(t *testing.T) {
	data := map[string]string{
		"id":    "test-user-2",
		"money": "1500.01",
	}

	err := testingTemplate("/increase", database.IncreaseAndDecrease, data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestDecreaseUser1(t *testing.T) {
	data := map[string]string{
		"id":    "test-user-1",
		"money": "15.01",
	}

	err := testingTemplate("/decrease", database.IncreaseAndDecrease, data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestDecreaseUser2(t *testing.T) {
	data := map[string]string{
		"id":    "test-user-2",
		"money": "20",
	}

	err := testingTemplate("/decrease", database.IncreaseAndDecrease, data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestRemittance(t *testing.T) {
	data := map[string]string{
		"from":  "test-user-1",
		"to":    "test-user-2",
		"money": "3.14",
	}

	err := testingTemplate("/remittance", database.Remittance, data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
