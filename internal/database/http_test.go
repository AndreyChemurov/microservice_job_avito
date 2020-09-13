package database_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"microservice_job_avito/internal/database"
	"net/http"
	"net/http/httptest"
	netURL "net/url"
	"os"
	"strings"
	"testing"

	"github.com/gchaincl/dotsql"
)

const (
	constURL string = "http://localhost:8000"
)

var (
	url string = ""
)

const (
	host     = "db"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

var postgresInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

type function func(http.ResponseWriter, *http.Request)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func prepareDatabase() error {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("start.sql")

	if err != nil {
		log.Fatal(err)
	}

	_, err = dot.Exec(db, "create-user-table")

	if err != nil {
		log.Fatal(err)
	}

	_, err = dot.Exec(db, "create-balance-table")

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func setup() {
	err := prepareDatabase()

	if err != nil {
		log.Fatal(err)
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

/* NORMAL BEHAVIOR TESTING */

func testingTemplateGET(route string, callingfn function, params map[string]string) error {
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

func testingTemplatePOST(route string, callingfn function, ct string, params map[string]string) error {
	url = (constURL + route)
	data := netURL.Values{}

	for k, v := range params {
		data.Set(k, v)
	}

	var (
		err     error
		request *http.Request
	)

	switch ct {
	case "application/json":
		j, _ := json.Marshal(params)
		request, err = http.NewRequest("POST", url, bytes.NewBuffer(j))

	case "application/x-www-form-urlencoded":
		request, err = http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	}

	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", ct)

	rec := httptest.NewRecorder()

	callingfn(rec, request)

	if want, got := 200, rec.Code; want != got {
		return fmt.Errorf("want %d, got %d; %v", want, got, rec.Body)
	}

	return nil
}

func TestGetBalance(t *testing.T) {
	data := map[string]string{
		"id": "test-user-1",
	}

	err := testingTemplateGET("/balance", database.GetBalance, data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGetBalanceCurrency(t *testing.T) {
	data := map[string]string{
		"id":       "test-user-1",
		"currency": "USD",
	}

	err := testingTemplateGET("/balance", database.GetBalance, data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGetBalanceCurrencyWWW(t *testing.T) {
	data := map[string]string{
		"id":       "test-user-1",
		"currency": "EUR",
	}

	err := testingTemplatePOST("/balance", database.GetBalance, "application/x-www-form-urlencoded", data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGetBalanceJSON(t *testing.T) {
	data := map[string]string{
		"id":       "test-user-1",
		"currency": "EUR",
	}

	err := testingTemplatePOST("/balance", database.GetBalance, "application/json", data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGetBalanceWWW(t *testing.T) {
	data := map[string]string{
		"id": "test-user-2",
	}

	err := testingTemplatePOST("/balance", database.GetBalance, "application/x-www-form-urlencoded", data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestIncreaseUser1(t *testing.T) {
	data := map[string]string{
		"id":    "test-user-1",
		"money": "1000.10",
	}

	err := testingTemplateGET("/increase", database.IncreaseAndDecrease, data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestIncreaseUser1JSON(t *testing.T) {
	data := map[string]string{
		"id":    "test-user-1",
		"money": "1106.87",
	}

	err := testingTemplatePOST("/increase", database.IncreaseAndDecrease, "application/json", data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestIncreaseUser1WWW(t *testing.T) {
	data := map[string]string{
		"id":    "test-user-2",
		"money": "2126.14",
	}

	err := testingTemplatePOST("/increase", database.IncreaseAndDecrease, "application/x-www-form-urlencoded", data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestIncreaseUser2(t *testing.T) {
	data := map[string]string{
		"id":    "test-user-2",
		"money": "1500.01",
	}

	err := testingTemplateGET("/increase", database.IncreaseAndDecrease, data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestDecreaseUser1(t *testing.T) {
	data := map[string]string{
		"id":    "test-user-1",
		"money": "15.01",
	}

	err := testingTemplateGET("/decrease", database.IncreaseAndDecrease, data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestDecreaseUser1JSON(t *testing.T) {
	data := map[string]string{
		"id":    "test-user-1",
		"money": "12",
	}

	err := testingTemplatePOST("/decrease", database.IncreaseAndDecrease, "application/json", data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestDecreaseUser2WWW(t *testing.T) {
	data := map[string]string{
		"id":    "test-user-2",
		"money": "15",
	}

	err := testingTemplatePOST("/decrease", database.IncreaseAndDecrease, "application/x-www-form-urlencoded", data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestDecreaseUser2(t *testing.T) {
	data := map[string]string{
		"id":    "test-user-2",
		"money": "20",
	}

	err := testingTemplateGET("/decrease", database.IncreaseAndDecrease, data)

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

	err := testingTemplateGET("/remittance", database.Remittance, data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestRemittanceJSON(t *testing.T) {
	data := map[string]string{
		"from":  "test-user-1",
		"to":    "test-user-2",
		"money": "2",
	}

	err := testingTemplatePOST("/remittance", database.Remittance, "application/json", data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestRemittanceWWW(t *testing.T) {
	data := map[string]string{
		"from":  "test-user-2",
		"to":    "test-user-1",
		"money": "1",
	}

	err := testingTemplatePOST("/remittance", database.Remittance, "application/x-www-form-urlencoded", data)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestNotFound(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8000/", nil)

	if err != nil {
		t.Fail()
	}

	rec := httptest.NewRecorder()

	database.NotFound(rec, request)

	if want, got := 404, rec.Code; want != got {
		t.Fail()
	}
}

/* ERROR TESTING */

func testingTemplateError(method string, path string, callingfn function, code int, data []byte, ct string) error {
	var (
		request *http.Request
		err     error
	)

	if data != nil {
		request, err = http.NewRequest(method, path, bytes.NewBuffer(data))
	} else {
		request, err = http.NewRequest(method, path, nil)
	}

	if ct != "" {
		request.Header.Add("Content-Type", ct)
	}

	if err != nil {
		return err
	}

	rec := httptest.NewRecorder()

	callingfn(rec, request)

	if got := rec.Code; code != got {
		return err
	}

	return nil
}

func TestGetBalanceUnknownUser(t *testing.T) {
	err := testingTemplateError("GET", "http://localhost:8000/balance?id=unknown-user-id", database.GetBalance, 500, nil, "")

	if err != nil {
		t.Fail()
	}
}

func TestGetBalanceWrongMethod(t *testing.T) {
	err := testingTemplateError("PUT", "http://localhost:8000/balance?id=unknown-user-id", database.GetBalance, 405, nil, "")

	if err != nil {
		t.Fail()
	}
}

func TestIncreaseWrongMethod(t *testing.T) {
	err := testingTemplateError("PUT", "http://localhost:8000/increase?id=unknown-user-id", database.IncreaseAndDecrease, 405, nil, "")

	if err != nil {
		t.Fail()
	}
}

func TestDecreaseWrongMethod(t *testing.T) {
	err := testingTemplateError("PUT", "http://localhost:8000/decrease?id=unknown-user-id", database.IncreaseAndDecrease, 405, nil, "")

	if err != nil {
		t.Fail()
	}
}

func TestRemittanceWrongMethod(t *testing.T) {
	err := testingTemplateError("PUT", "http://localhost:8000/remittance?from=test-user-1&to=test-user-2&money=1", database.Remittance, 405, nil, "")

	if err != nil {
		t.Fail()
	}
}

func TestDecreaseMore(t *testing.T) {
	err := testingTemplateError("GET", "http://localhost:8000/decrease?id=test-user-1&money=100000", database.IncreaseAndDecrease, 500, nil, "")

	if err != nil {
		t.Fail()
	}
}

func TestPOSTCTNoURL(t *testing.T) {
	err := testingTemplateError("POST", "http://localhost:8000/balance?id=test-user-1", database.GetBalance, 400, nil, "application/json")

	if err != nil {
		t.Fail()
	}
}

func TestGETCTNoBody(t *testing.T) {
	err := testingTemplateError("GET", "http://localhost:8000/balance", database.GetBalance, 405, nil, "application/json")

	if err != nil {
		t.Fail()
	}
}

func TestPOSTCTNoURLIncrease(t *testing.T) {
	err := testingTemplateError("POST", "http://localhost:8000/increase?id=test-user-1&money=1", database.IncreaseAndDecrease, 400, nil, "application/json")

	if err != nil {
		t.Fail()
	}
}

func TestGETCTNoBodyIncrease(t *testing.T) {
	err := testingTemplateError("GET", "http://localhost:8000/increase", database.IncreaseAndDecrease, 405, nil, "application/json")

	if err != nil {
		t.Fail()
	}
}

func TestPOSTCTNoURLDecrease(t *testing.T) {
	err := testingTemplateError("POST", "http://localhost:8000/decrease?id=test-user-1&money=1", database.IncreaseAndDecrease, 400, nil, "application/json")

	if err != nil {
		t.Fail()
	}
}

func TestGETCTNoBodyDecrease(t *testing.T) {
	err := testingTemplateError("GET", "http://localhost:8000/decrease", database.IncreaseAndDecrease, 405, nil, "application/json")

	if err != nil {
		t.Fail()
	}
}

func TestPOSTCTNoURLRemittance(t *testing.T) {
	err := testingTemplateError("POST", "http://localhost:8000/remittance?from=test-user-1&to=test-user-2&money=1", database.Remittance, 400, nil, "application/json")

	if err != nil {
		t.Fail()
	}
}

func TestGETCTNoBodyRemmitance(t *testing.T) {
	err := testingTemplateError("GET", "http://localhost:8000/remittance", database.Remittance, 405, nil, "application/json")

	if err != nil {
		t.Fail()
	}
}

func TestGetBalanceWrongParams(t *testing.T) {
	err := testingTemplateError("GET", "http://localhost:8000/balance?unknown=param", database.GetBalance, 400, nil, "")

	if err != nil {
		t.Fail()
	}
}

func TestIncreaseWrongParams(t *testing.T) {
	err := testingTemplateError("GET", "http://localhost:8000/increase?unknown=param", database.IncreaseAndDecrease, 400, nil, "")

	if err != nil {
		t.Fail()
	}
}

func TestRemittanceWrongParams(t *testing.T) {
	err := testingTemplateError("GET", "http://localhost:8000/remittance?unknown=param", database.Remittance, 400, nil, "")

	if err != nil {
		t.Fail()
	}
}

func TestIncreaseWrongParams2(t *testing.T) {
	err := testingTemplateError("GET", "http://localhost:8000/increase?id=test-user-1&money=unknown", database.IncreaseAndDecrease, 400, nil, "")

	if err != nil {
		t.Fail()
	}
}

func TestRemittanceWrongParams2(t *testing.T) {
	err := testingTemplateError("GET", "http://localhost:8000/remittance?from=test-user-1&to=test-user-2&money=unknown", database.Remittance, 400, nil, "")

	if err != nil {
		t.Fail()
	}
}

func TestUnknownCT(t *testing.T) {
	err := testingTemplateError("POST", "http://localhost:8000/balance", database.GetBalance, 400, []byte(`{"id":"test-user-1"}`), "application/unknown")

	if err != nil {
		t.Fail()
	}
}

func TestGetBalanceWrongJSON(t *testing.T) {
	err := testingTemplateError("POST", "http://localhost:8000/balance", database.GetBalance, 400, []byte(`{id:test-user-1}`), "application/json")

	if err != nil {
		t.Fail()
	}
}

func TestIncreaseWrongJSON(t *testing.T) {
	err := testingTemplateError("POST", "http://localhost:8000/increase", database.IncreaseAndDecrease, 400, []byte(`{id:test-user-1}`), "application/json")

	if err != nil {
		t.Fail()
	}
}

func TestRemittanceWrongJSON(t *testing.T) {
	err := testingTemplateError("POST", "http://localhost:8000/remittance", database.Remittance, 400, []byte(`{id:test-user-1}`), "application/json")

	if err != nil {
		t.Fail()
	}
}

func TestRemittanceUnknownUserFrom(t *testing.T) {
	err := testingTemplateError("POST", "http://localhost:8000/remittance", database.Remittance, 500, []byte(`{"from":"unknown-test-user-from","to":"test-user-1","money":"1"}`), "application/json")

	if err != nil {
		t.Fail()
	}
}

func TestRemittanceUnknownUserTo(t *testing.T) {
	err := testingTemplateError("POST", "http://localhost:8000/remittance", database.Remittance, 500, []byte(`{"from":"test-user-1","to":"unknown-test-user-1","money":"1"}`), "application/json")

	if err != nil {
		t.Fail()
	}
}

func TestWrongCurrency(t *testing.T) {
	err := testingTemplateError("POST", "http://localhost:8000/balance", database.GetBalance, 500, []byte(`{"id":"test-user-1","currency":"UNKNOWN"}`), "application/json")

	if err != nil {
		t.Fail()
	}
}
