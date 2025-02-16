package pgxtxpool

import "github.com/google/uuid"

// generateID will generate unique id for transaction ID
func generateID() TxID {
	return TxID(uuid.New().String())
}
