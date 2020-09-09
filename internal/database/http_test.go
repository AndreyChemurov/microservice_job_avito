package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gchaincl/dotsql"
	_ "github.com/lib/pq" //
)

var (
	// recorder *httptest.ResponseRecorder
	request  *http.Request
	response *http.Response
	client   *http.Client

	url string
)

const (
	constURL string = "http://localhost:8000"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("../../start.sql")

	if err != nil {
		panic(err)
	}

	_, err = dot.Exec(db, "create-user", "avito_test_user_1")

	if err != nil {
		panic(err)
	}

	_, err = dot.Exec(db, "create-user", "avito_test_user_2")

	if err != nil {
		panic(err)
	}

	_, err = dot.Exec(db, "create-balance", "avito_test_user_1")

	if err != nil {
		panic(err)
	}

	_, err = dot.Exec(db, "create-balance", "avito_test_user_2")

	if err != nil {
		panic(err)
	}
}

func teardown() {
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	dot, err := dotsql.LoadFromFile("../../start.sql")

	_, err = dot.Exec(db, "drop-balance", "avito_test_user_1")

	if err != nil {
		panic(err)
	}

	_, err = dot.Exec(db, "drop-balance", "avito_test_user_2")

	if err != nil {
		panic(err)
	}

	_, err = dot.Exec(db, "drop-user", "avito_test_user_1")

	if err != nil {
		panic(err)
	}

	_, err = dot.Exec(db, "drop-user", "avito_test_user_2")

	if err != nil {
		panic(err)
	}
}

func testingTemplate(t *testing.T, route string, params map[string]string) {
	url += (route + "?")

	for k, v := range params {
		url += (k + "=" + v + "&")
	}

	url = url[:len(url)-1]

	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		t.Errorf("No Database Connection")
	}

	defer db.Close()

	// jsonData := map[string]string{"id": "someuser"}
	// jsonValue, _ := json.Marshal(jsonData)

	// recorder := httptest.NewRecorder()

	// log.Println(url)

	request, _ = http.NewRequest("GET", url, nil) // , bytes.NewBuffer(jsonValue)

	url = constURL

	// request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client = &http.Client{}
	response, err = client.Do(request)

	if err != nil {
		fmt.Println("HTTP call failed:", err)
		t.Fail()
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	bodyStr := string(body)

	if response.StatusCode != http.StatusOK {
		fmt.Println("Non-OK HTTP status:", response.StatusCode, bodyStr)
		t.Fail()
		return
	}
}

func TestGetBalance(t *testing.T) {
	data := map[string]string{
		"id":       "avito_test_user_1",
		"currency": "USD",
	}

	testingTemplate(t, "/balance", data)
}

func TestIncrease(t *testing.T) {
	data := map[string]string{
		"id":    "avito_test_user_1",
		"money": "1000",
	}

	testingTemplate(t, "/increase", data)
}

func TestDecrease(t *testing.T) {
	data := map[string]string{
		"id":    "avito_test_user_1",
		"money": "100",
	}

	testingTemplate(t, "/decrease", data)
}

func TestRemittance(t *testing.T) {
	data := map[string]string{
		"from":  "avito_test_user_1",
		"to":    "avito_test_user_2",
		"money": "10",
	}

	testingTemplate(t, "/remittance", data)
}
