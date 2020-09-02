package database

type rm map[string]map[string]string // rm stands for Response Message

// var name: <rm_name><http_status_code>rm

// NotFound404rm ...
var NotFound404rm rm = rm{
	"error": {
		"status_code":    "404",
		"status_message": "Not Found",
	},
}

// NoDatabaseConnection500rm ...
var NoDatabaseConnection500rm rm = rm{
	"error": {
		"status_code":    "500",
		"status_message": "Service cannot connect to database",
	},
}

// UserDoesNotExist400rm ...
var UserDoesNotExist400rm rm = rm{
	"error": {
		"status_code":    "400",
		"status_message": "User does not exist",
	},
}

// WrongParams400rm ...
var WrongParams400rm rm = rm{
	"error": {
		"status_code":    "400",
		"status_message": "Wrong parameter(s)",
	},
}

// InternalServerError500rm ...
var InternalServerError500rm rm = rm{
	"error": {
		"status_code":    "500",
		"status_message": "Internal Server Error",
	},
}
