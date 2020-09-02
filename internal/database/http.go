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

func getBalance(w http.ResponseWriter, r *http.Request) {
	userIDFromRequest := r.URL.Query().Get("id")

	if userIDFromRequest == "" {
		js, _ := json.Marshal(WrongParams400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	balance, status := _getBalance(userIDFromRequest)

	w.WriteHeader(status)
	w.Write(balance)
	return
}

func increase(w http.ResponseWriter, r *http.Request) {
	userIDFromRequest := r.URL.Query().Get("id")
	moneyFromRequest := r.URL.Query().Get("money")

	if userIDFromRequest == "" || moneyFromRequest == "" {
		js, _ := json.Marshal(WrongParams400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	money, err := strconv.ParseFloat(moneyFromRequest, 64)

	if err != nil {
		js, _ := json.Marshal(WrongParams400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	userID, status := _increase(userIDFromRequest, math.Round(money*100)/100)

	w.WriteHeader(status)
	w.Write(userID)
	return
}

func decrease(w http.ResponseWriter, r *http.Request) {
	//
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
