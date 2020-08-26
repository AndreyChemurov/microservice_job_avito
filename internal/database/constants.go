package database

type rm map[string]map[string]string // rm stands for Response Message

// var name: <http_status_name><http_status_code>rm

// NotFound404rm ...
var NotFound404rm rm = rm{
	"error": {
		"status_code":    "404",
		"status_message": "Not Found",
	},
}
