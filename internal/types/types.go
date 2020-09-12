package types

/* Response Types */

// Balance - применяется в  GetBalance, IncreaseAndDecrease
type Balance struct {
	Balance float64
	Status  int
}

// Remittance применяется в Remittance
type Remittance struct {
	BalanceFrom float64
	BalanceTo   float64
	Status      int
}

/* Request Types */

// UserIDBalance - запрос в /balance
type UserIDBalance struct {
	ID       string `json:"id"`
	Currency string `json:"currency"`
}

// IncreaseDecrease - запросы в /increase и /decrease
type IncreaseDecrease struct {
	ID    string `json:"id"`
	Money string `json:"money"`
}

// RemittanceRequest - запрос в /remittance
type RemittanceRequest struct {
	IDFrom string `json:"from"`
	IDTo   string `json:"to"`
	Money  string `json:"money"`
}

// Currency - запрос в /balance с параметром "currency"
type Currency struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}
