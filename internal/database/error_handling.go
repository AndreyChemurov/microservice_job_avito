package database

import (
	"encoding/json"
	"strconv"
)

type rm map[string]map[string]string // rm stands for Response Message

// ErrorType оборачивает данные (http статус код + сообщение) об ошибке
// в json и возвращает его вызывающему хэндлеру
func ErrorType(status int, message string) (response []byte) {
	err := rm{
		"error": {
			"status_code":    strconv.Itoa(status),
			"status_message": message,
		},
	}

	returnJSON, _ := json.Marshal(err)
	return returnJSON
}
