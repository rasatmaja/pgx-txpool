package pgxtxpool

// TxID an identifier for a transaction
type TxID string

// TxContextID an indentifier for a transaction ID in context
type TxContextID string

// ContextTxKey a key for a transaction ID in context
const ContextTxKey TxContextID = "TX_POLL_ID"
