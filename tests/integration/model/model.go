package model

// User is a struct that represent a user
type User struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Balance       float64 `json:"balance"`
	BalanceChange float64 `json:"balance_change"`
}

// Transaction is a struct that represent a transaction
type Transaction struct {
	ID     string  `json:"id"`
	UserID string  `json:"user_id"`
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}

// TransactionTransfer is a struct that represent a transaction transfer
type TransactionTransfer struct {
	ID                       string  `json:"id"`
	TransactionOriginID      string  `json:"transaction_origin_id"`
	TransactionDestinationID string  `json:"transaction_destination_id"`
	Amount                   float64 `json:"amount"`
}
