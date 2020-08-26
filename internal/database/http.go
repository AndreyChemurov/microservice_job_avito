package database

import (
	"encoding/json"
	"log"
	"net/http"
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
	// TODO:
	//	1. How to send data: url params or post body? (both?)
	//	2. increase + decrease == one function
	//	3. Constants + exceptions
	//	4. Parse right data type from request to avoid sql injection

	// EXTRA:
	//	1. Other currencies (currency param)
}

func remittance(w http.ResponseWriter, r *http.Request) {
	//
}

func increase(w http.ResponseWriter, r *http.Request) {
	//
}

func decrease(w http.ResponseWriter, r *http.Request) {
	//
}

func notFound(w http.ResponseWriter, r *http.Request) {
	js, _ := json.Marshal(NotFound404rm)
	w.WriteHeader(http.StatusNotFound)
	w.Write(js)
	return
}
