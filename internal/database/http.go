package database

import (
	"encoding/json"
	"log"
	"math"
	"microservice_job_avito/internal/types"
	"net/http"
	"strconv"
	"strings"
)

// PathHandler ...
func PathHandler() {
	http.HandleFunc("/balance", GetBalance)
	http.HandleFunc("/remittance", Remittance)
	http.HandleFunc("/increase", IncreaseAndDecrease)
	http.HandleFunc("/decrease", IncreaseAndDecrease)

	http.HandleFunc("/", NotFound)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

// GetBalance - Метод получения текущего баланса пользователя.
// Аргументы:
//		id: уникальный идентификатор пользователя;
// Возвращаемые значения:
//		balance: текущий баланс пользователя;
func GetBalance(w http.ResponseWriter, r *http.Request) {
	var (
		contentType string = r.Header.Get("Content-Type")

		userIDFromRequest   string
		currencyFromRequest string
		currencyFlag        bool = false
		js                  []byte

		balance []byte
		status  int
	)

	if r.Method != "POST" && r.Method != "GET" {
		js, _ = json.Marshal(MethodNotAllowed405rm)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(js)
		return
	}

	if contentType != "" {
		// Метод POST, Content-Type, но URL-параметры вместо тела
		if _, exist := r.URL.Query()["id"]; exist {
			js, _ = json.Marshal(URLParams400rm)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(js)
			return
		}

		// Метод GET и пустое тело
		if r.Method == "GET" && userIDFromRequest == "" {
			js, _ = json.Marshal(ContentTypeGETMethod405rm)
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(js)
			return
		}
	}

	// Парсинг параметров
	if contentType == "application/x-www-form-urlencoded" {
		r.ParseForm()
		userIDFromRequest = r.Form.Get("id")
		currencyFromRequest = r.Form.Get("currency")

		if currencyFromRequest != "" {
			currencyFlag = true
		}

	} else if contentType == "application/json" {
		decoder := json.NewDecoder(r.Body)

		var reqData types.UserIDBalance

		if err := decoder.Decode(&reqData); err != nil {
			js, _ = json.Marshal(BadJSON400rm)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(js)
			return
		}

		userIDFromRequest = reqData.ID
		currencyFromRequest = reqData.Currency

		if currencyFromRequest != "" {
			currencyFlag = true
		}

	} else if contentType == "" {
		userIDFromRequest = r.URL.Query().Get("id")
		currencyFromRequest = r.URL.Query().Get("currency")

		if currencyFromRequest != "" {
			currencyFlag = true
		}

	} else {
		js, _ = json.Marshal(UnknownContentType400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	// Если параметры введены неверно
	if userIDFromRequest == "" {
		js, _ = json.Marshal(WrongParams400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	balance, status = _getBalance(userIDFromRequest, currencyFlag, currencyFromRequest)

	w.WriteHeader(status)
	w.Write(balance)
	return
}

// IncreaseAndDecrease - Метод начисления и списания средств.
// Аргументы:
//		id: уникальный идентификатор пользователя;
//		money: количество средст для зачисления/списания;
// Возвращаемые значения:
//		balance: текущий баланс пользователя;
func IncreaseAndDecrease(w http.ResponseWriter, r *http.Request) {
	var (
		contentType string = r.Header.Get("Content-Type")

		userIDFromRequest string
		moneyFromRequest  string
		js                []byte

		balance []byte
		status  int
	)

	if r.Method != "POST" && r.Method != "GET" {
		js, _ = json.Marshal(MethodNotAllowed405rm)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(js)
		return
	}

	params := [2]string{"id", "money"}

	if contentType != "" {
		// Метод POST, Content-Type, но URL-параметры вместо тела
		for _, param := range params {
			if _, exist := r.URL.Query()[param]; exist {
				js, _ = json.Marshal(URLParams400rm)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(js)
				return
			}
		}

		// Метод GET и пустое тело
		if r.Method == "GET" && (userIDFromRequest == "" || moneyFromRequest == "") {
			js, _ = json.Marshal(ContentTypeGETMethod405rm)
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(js)
			return
		}
	}

	// Парсинг параметров
	if contentType == "application/x-www-form-urlencoded" {
		r.ParseForm()

		userIDFromRequest = r.Form.Get("id")
		moneyFromRequest = r.Form.Get("money")

	} else if contentType == "application/json" {
		decoder := json.NewDecoder(r.Body)

		var reqData types.IncreaseDecrease

		if err := decoder.Decode(&reqData); err != nil {
			js, _ = json.Marshal(BadJSON400rm)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(js)
			return
		}

		userIDFromRequest = reqData.ID
		moneyFromRequest = reqData.Money

	} else if contentType == "" {
		userIDFromRequest = r.URL.Query().Get("id")
		moneyFromRequest = r.URL.Query().Get("money")
	} else {
		js, _ = json.Marshal(UnknownContentType400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	// Если параметры введены неверно
	if userIDFromRequest == "" || moneyFromRequest == "" {
		js, _ := json.Marshal(WrongParams400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	// Если в указанной сумме используется запятая, а не точка
	moneyFromRequest = strings.Replace(moneyFromRequest, ",", ".", -1)
	money, err := strconv.ParseFloat(moneyFromRequest, 64)

	if err != nil || money < 0 {
		js, _ := json.Marshal(WrongParams400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	if r.URL.Path == "/increase" {
		balance, status = _increase(userIDFromRequest, math.Round(money*100)/100)
	} else if r.URL.Path == "/decrease" {
		balance, status = _decrease(userIDFromRequest, math.Round(money*100)/100)
	}

	w.WriteHeader(status)
	w.Write(balance)
	return
}

// Remittance - Метод перевода средств от пользователя к пользователю.
// Аргументы:
//		from: уникальный идентификатор пользователя, с баланса которого надо списать средства;
//		to: уникальный идентификатор пользователя, на баланс которого надо перечислить средства;
//		money: количество средст для зачисления;
// Возвращаемые значения:
//		balance: текущие балансы обоих пользователей;
func Remittance(w http.ResponseWriter, r *http.Request) {
	var (
		contentType string = r.Header.Get("Content-Type")

		userFromIDRequest string
		userToIDRequest   string
		moneyFromRequest  string
		js                []byte

		balance []byte
		status  int
	)

	if r.Method != "POST" && r.Method != "GET" {
		js, _ = json.Marshal(MethodNotAllowed405rm)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(js)
		return
	}

	params := [3]string{"from", "to", "money"}

	if contentType != "" {
		// Метод POST, Content-Type, но URL-параметры вместо тела
		for _, param := range params {
			if _, exist := r.URL.Query()[param]; exist {
				js, _ = json.Marshal(URLParams400rm)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(js)
				return
			}
		}

		// Метод GET и пустое тело
		if r.Method == "GET" && (userFromIDRequest == "" || userToIDRequest == "" || moneyFromRequest == "") {
			js, _ = json.Marshal(ContentTypeGETMethod405rm)
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(js)
			return
		}
	}

	// Парсинг параметров
	if contentType == "application/x-www-form-urlencoded" {
		r.ParseForm()

		userFromIDRequest = r.Form.Get("from")
		userToIDRequest = r.Form.Get("to")
		moneyFromRequest = r.Form.Get("money")

	} else if contentType == "application/json" {
		decoder := json.NewDecoder(r.Body)

		var reqData types.RemittanceRequest

		if err := decoder.Decode(&reqData); err != nil {
			js, _ = json.Marshal(BadJSON400rm)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(js)
			return
		}

		userFromIDRequest = reqData.IDFrom
		userToIDRequest = reqData.IDTo
		moneyFromRequest = reqData.Money

	} else if contentType == "" {
		userFromIDRequest = r.URL.Query().Get("from")
		userToIDRequest = r.URL.Query().Get("to")
		moneyFromRequest = r.URL.Query().Get("money")
	} else {
		js, _ = json.Marshal(UnknownContentType400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	// Если параметры введены неверно
	if userFromIDRequest == "" || userToIDRequest == "" || moneyFromRequest == "" {
		js, _ := json.Marshal(WrongParams400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	// Если в указанной сумме используется запятая, а не точка
	moneyFromRequest = strings.Replace(moneyFromRequest, ",", ".", -1)
	money, err := strconv.ParseFloat(moneyFromRequest, 64)

	if err != nil || money < 0 {
		js, _ := json.Marshal(WrongParams400rm)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	balance, status = _remittance(userFromIDRequest, userToIDRequest, math.Round(money*100)/100)

	w.WriteHeader(status)
	w.Write(balance)
	return
}

// NotFound - Метод, который вызывается, если указанный путь не найден.
func NotFound(w http.ResponseWriter, r *http.Request) {
	js, _ := json.Marshal(NotFound404rm)
	w.WriteHeader(http.StatusNotFound)
	w.Write(js)
	return
}
