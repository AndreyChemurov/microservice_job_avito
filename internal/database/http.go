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
		responseJSON        []byte

		balance *types.Balance
		status  int
	)

	if r.Method != "POST" && r.Method != "GET" {
		responseJSON = ErrorType(405, "Method not allowed: use GET or POST")

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(responseJSON)

		return
	}

	if contentType != "" {
		// Метод POST, Content-Type, но URL-параметры вместо тела
		if _, exist := r.URL.Query()["id"]; exist {
			responseJSON = ErrorType(400, "Need body parameters, not URL")

			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseJSON)

			return
		}

		// Метод GET и пустое тело
		if r.Method == "GET" && userIDFromRequest == "" {
			responseJSON = ErrorType(405, "Method not allowed: use POST-method and body parameters with Content-Type")

			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(responseJSON)

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
			responseJSON = ErrorType(400, "Invalid JSON format")

			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseJSON)

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
		responseJSON = ErrorType(400, "Unknown Content-Type")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	// Если параметры введены неверно
	if userIDFromRequest == "" {
		responseJSON = ErrorType(400, "Wrong parameter(s)")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	balance, status, err := getBalance(userIDFromRequest, currencyFlag, currencyFromRequest)

	if err != nil {
		responseJSON = ErrorType(status, err.Error())

		w.WriteHeader(status)
		w.Write(responseJSON)

		return
	}

	responseJSON, _ = json.Marshal(balance)

	w.WriteHeader(200)
	w.Write(responseJSON)

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
		responseJSON      []byte

		balance *types.Balance
		status  int
	)

	if r.Method != "POST" && r.Method != "GET" {
		responseJSON = ErrorType(405, "Method not allowed: use GET or POST")

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(responseJSON)

		return
	}

	params := [2]string{"id", "money"}

	if contentType != "" {
		// Метод POST, Content-Type, но URL-параметры вместо тела
		for _, param := range params {
			if _, exist := r.URL.Query()[param]; exist {
				responseJSON = ErrorType(400, "Need body parameters, not URL")

				w.WriteHeader(http.StatusBadRequest)
				w.Write(responseJSON)

				return
			}
		}

		// Метод GET и пустое тело
		if r.Method == "GET" && (userIDFromRequest == "" || moneyFromRequest == "") {
			responseJSON = ErrorType(405, "Method not allowed: use POST-method and body parameters with Content-Type")

			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(responseJSON)

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
			responseJSON = ErrorType(400, "Invalid JSON format")

			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseJSON)

			return
		}

		userIDFromRequest = reqData.ID
		moneyFromRequest = reqData.Money

	} else if contentType == "" {
		userIDFromRequest = r.URL.Query().Get("id")
		moneyFromRequest = r.URL.Query().Get("money")
	} else {
		responseJSON = ErrorType(400, "Unknown Content-Type")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	// Если параметры введены неверно
	if userIDFromRequest == "" || moneyFromRequest == "" {
		responseJSON = ErrorType(400, "Wrong parameter(s)")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	// Если в указанной сумме используется запятая, а не точка
	moneyFromRequest = strings.Replace(moneyFromRequest, ",", ".", -1)
	money, err := strconv.ParseFloat(moneyFromRequest, 64)

	if err != nil || money < 0 {
		responseJSON = ErrorType(400, "Wrong parameter(s)")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	if r.URL.Path == "/increase" {
		balance, status, err = increase(userIDFromRequest, math.Round(money*100)/100)

	} else if r.URL.Path == "/decrease" {
		balance, status, err = decrease(userIDFromRequest, math.Round(money*100)/100)
	}

	if err != nil {
		responseJSON = ErrorType(status, err.Error())

		w.WriteHeader(status)
		w.Write(responseJSON)

		return
	}

	responseJSON, _ = json.Marshal(balance)

	w.WriteHeader(status)
	w.Write(responseJSON)

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
		responseJSON      []byte

		balance *types.Remittance
		status  int
	)

	if r.Method != "POST" && r.Method != "GET" {
		responseJSON = ErrorType(405, "Method not allowed: use GET or POST")

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(responseJSON)

		return
	}

	params := [3]string{"from", "to", "money"}

	if contentType != "" {
		// Метод POST, Content-Type, но URL-параметры вместо тела
		for _, param := range params {
			if _, exist := r.URL.Query()[param]; exist {
				responseJSON = ErrorType(400, "Need body parameters, not URL")

				w.WriteHeader(http.StatusBadRequest)
				w.Write(responseJSON)

				return
			}
		}

		// Метод GET и пустое тело
		if r.Method == "GET" && (userFromIDRequest == "" || userToIDRequest == "" || moneyFromRequest == "") {
			responseJSON = ErrorType(405, "Method not allowed: use POST-method and body parameters with Content-Type")

			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(responseJSON)

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
			responseJSON = ErrorType(400, "Invalid JSON format")

			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseJSON)

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
		responseJSON = ErrorType(400, "Unknown Content-Type")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	// Если параметры введены неверно
	if userFromIDRequest == "" || userToIDRequest == "" || moneyFromRequest == "" {
		responseJSON = ErrorType(400, "Wrong parameter(s)")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	// Если в указанной сумме используется запятая, а не точка
	moneyFromRequest = strings.Replace(moneyFromRequest, ",", ".", -1)
	money, err := strconv.ParseFloat(moneyFromRequest, 64)

	if err != nil || money < 0 {
		responseJSON = ErrorType(400, "Wrong parameter(s)")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	balance, status, err = remittance(userFromIDRequest, userToIDRequest, math.Round(money*100)/100)

	if err != nil {
		responseJSON = ErrorType(status, err.Error())

		w.WriteHeader(status)
		w.Write(responseJSON)

		return
	}

	responseJSON, _ = json.Marshal(balance)

	w.WriteHeader(status)
	w.Write(responseJSON)

	return
}

// NotFound - Метод, который вызывается, если указанный путь не найден.
func NotFound(w http.ResponseWriter, r *http.Request) {
	responseJSON := ErrorType(404, "Not found")

	w.WriteHeader(http.StatusNotFound)
	w.Write(responseJSON)

	return
}
