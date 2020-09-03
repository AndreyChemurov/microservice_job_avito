package database

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// PathHandler ...
func PathHandler() {
	http.HandleFunc("/balance", getBalance)
	http.HandleFunc("/remittance", remittance)
	http.HandleFunc("/increase", increaseAndDecrease)
	http.HandleFunc("/decrease", increaseAndDecrease)

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

func increaseAndDecrease(w http.ResponseWriter, r *http.Request) {
	userIDFromRequest := r.URL.Query().Get("id")
	moneyFromRequest := r.URL.Query().Get("money")

	if userIDFromRequest == "" || moneyFromRequest == "" {
		js, _ := json.Marshal(WrongParams400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	moneyFromRequest = strings.Replace(moneyFromRequest, ",", ".", -1)
	money, err := strconv.ParseFloat(moneyFromRequest, 64)

	if err != nil {
		js, _ := json.Marshal(WrongParams400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	balance := make([]byte, 0)
	status := 0

	if r.URL.Path == "/increase" {
		balance, status = _increase(userIDFromRequest, math.Round(money*100)/100)
	} else if r.URL.Path == "/decrease" {
		balance, status = _decrease(userIDFromRequest, math.Round(money*100)/100)
	}

	w.WriteHeader(status)
	w.Write(balance)
	return
}

func remittance(w http.ResponseWriter, r *http.Request) {
	userFromIDRequest := r.URL.Query().Get("from")
	userToIDRequest := r.URL.Query().Get("to")
	moneyFromRequest := r.URL.Query().Get("money")

	if userFromIDRequest == "" || userToIDRequest == "" || moneyFromRequest == "" {
		js, _ := json.Marshal(WrongParams400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	moneyFromRequest = strings.Replace(moneyFromRequest, ",", ".", -1)
	money, err := strconv.ParseFloat(moneyFromRequest, 64)

	if err != nil {
		js, _ := json.Marshal(WrongParams400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	balance := make([]byte, 0)
	status := 0

	balance, status = _remittance(userFromIDRequest, userToIDRequest, math.Round(money*100)/100)

	w.WriteHeader(status)
	w.Write(balance)
	return
}

func notFound(w http.ResponseWriter, r *http.Request) {
	js, _ := json.Marshal(NotFound404rm)
	w.WriteHeader(http.StatusNotFound)
	w.Write(js)
	return
}
