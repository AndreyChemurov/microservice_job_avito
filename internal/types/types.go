package types

/* Response Types */

// Balance ...
type Balance struct {
	Bal float64
}

// Remittance ...
type Remittance struct {
	BalanceFrom float64
	BalanceTo   float64
}

/* Request Types */

// UserIDBalance ...
type UserIDBalance struct {
	ID string `json:"id"`
}

// IncreaseDecrease ...
type IncreaseDecrease struct {
	ID    string `json:"id"`
	Money string `json:"money"`
}

// RemittanceRequest ...
type RemittanceRequest struct {
	IDFrom string `json:"from"`
	IDTo   string `json:"to"`
	Money  string `json:"money"`
}

// Currency ...
type Currency struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}
