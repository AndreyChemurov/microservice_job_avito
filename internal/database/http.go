package database

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
)

// PathHandler ...
func PathHandler() {
	http.HandleFunc("/balance", getBalance)
	http.HandleFunc("/remittance", remittance)
	http.HandleFunc("/increase", increase)
	http.HandleFunc("/decrease", decrease)

	http.HandleFunc("/", notFound)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

// TODO:
//	2. increase + decrease == one function
//	3. Constants + exceptions
//	4. Parse right data type from request to avoid sql injection
//	5. sync-coming decrement request, prevent negative value

// EXTRA:
//	1. Other currencies (currency param)

func getBalance(w http.ResponseWriter, r *http.Request) {
	userIDFromRequest := r.URL.Query().Get("id")

	userID, err := strconv.ParseUint(userIDFromRequest, 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	_ = _getBalance(userID)
}

func increase(w http.ResponseWriter, r *http.Request) {
	userIDFromRequest := r.URL.Query().Get("id")
	moneyFromRequest := r.URL.Query().Get("money")

	userID, err := strconv.ParseUint(userIDFromRequest, 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	m, err := strconv.ParseFloat(moneyFromRequest, 64)

	if err != nil {
		log.Fatal(err)
	}

	money := math.Round(m*100) / 100

	_increase(userID, money) //
}

func decrease(w http.ResponseWriter, r *http.Request) {
	userIDFromRequest := r.URL.Query().Get("id")
	moneyFromRequest := r.URL.Query().Get("money")

	userID, err := strconv.ParseUint(userIDFromRequest, 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	m, err := strconv.ParseFloat(moneyFromRequest, 64)

	if err != nil {
		log.Fatal(err)
	}

	money := math.Round(m*100) / 100

	_decrease(userID, money)
}

func remittance(w http.ResponseWriter, r *http.Request) {
	//
}

func notFound(w http.ResponseWriter, r *http.Request) {
	js, _ := json.Marshal(NotFound404rm)
	w.WriteHeader(http.StatusNotFound)
	w.Write(js)
	return
}
