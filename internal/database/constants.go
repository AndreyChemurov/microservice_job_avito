package database

type rm map[string]map[string]interface{} // rm stands for Response Message

// var name: <rm_name><http_status_code>rm

// NotFound404rm ...
var NotFound404rm rm = rm{
	"error": {
		"status_code":    404,
		"status_message": "Not Found",
	},
}

// NoDatabaseConnection500rm ...
var NoDatabaseConnection500rm rm = rm{
	"error": {
		"status_code":    500,
		"status_message": "Service cannot connect to database",
	},
}

// UserDoesNotExist400rm ...
var UserDoesNotExist400rm rm = rm{
	"error": {
		"status_code":    400,
		"status_message": "User does not exist",
	},
}

// WrongParams400rm ...
var WrongParams400rm rm = rm{
	"error": {
		"status_code":    400,
		"status_message": "Wrong parameter(s)",
	},
}

// InternalServerError500rm ...
var InternalServerError500rm rm = rm{
	"error": {
		"status_code":    500,
		"status_message": "Internal Server Error",
	},
}

// DecreaseMore400rm ...
var DecreaseMore400rm rm = rm{
	"error": {
		"status_code":    400,
		"status_message": "Operation not available: Debit exceeds the balance",
	},
}

// MethodNotAllowed405rm ...
var MethodNotAllowed405rm rm = rm{
	"error": {
		"status_code":    405,
		"status_message": "Method not allowed: use GET or POST",
	},
}

// ContentTypeGETMethod405rm ...
var ContentTypeGETMethod405rm rm = rm{
	"error": {
		"status_code":    405,
		"status_message": "Method not allowed: use POST-method and body parameters with Content-Type",
	},
}

// URLParams400rm ...
var URLParams400rm rm = rm{
	"error": {
		"status_code":    400,
		"status_message": "POST-method and Content-Type require body parameters, not URL",
	},
}

// BadJSON400rm ...
var BadJSON400rm rm = rm{
	"error": {
		"status_code":    400,
		"status_message": "Invalid JSON format",
	},
}

// UnknownContentType400rm ...
var UnknownContentType400rm rm = rm{
	"error": {
		"status_code":    400,
		"status_message": "Unknown Content-Type",
	},
}

// WrongBaseFormat400rm ...
var WrongBaseFormat400rm rm = rm{
	"error": {
		"status_code":    400,
		"status_message": "Wrong base format for currency, example: USD, EUR, RUB",
	},
}
