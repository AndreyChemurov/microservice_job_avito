package currency

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"microservice_job_avito/internal/types"
	"net/http"
)

// Currency - парсер курса валюты через сервис exchangeratesapi
func Currency(balance float64, base string) (float64, error) {
	currencyReq := fmt.Sprintf("https://api.exchangeratesapi.io/latest?base=%s", base)

	resp, err := http.Get(currencyReq)

	if err != nil {
		return balance, errors.New("No response from exchangeratesapi")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return balance, errors.New("Service can't read exchangeratesapi's response body")
	}

	var c types.Currency
	err = json.Unmarshal(body, &c)

	if err != nil {
		return balance, errors.New("Currency's JSON-Unmarshal problem")
	}

	if len(c.Rates) == 0 {
		return balance, errors.New("Currency's JSON-Unmarshal problem")
	}

	balance /= c.Rates["RUB"]
	balance = math.Round(balance*100) / 100

	return balance, nil
}
