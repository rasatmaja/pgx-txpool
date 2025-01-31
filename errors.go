package pgxtxpool

import "fmt"

// ErrTxPool will indicate a pgx tx pool error
// this is a parent error (L1)
var ErrTxPool = fmt.Errorf("pgx tx pool error")

// ErrTxPoolIDNotFound will indicate that current context does not have transaction id with key ContextTxKey
// this is a child error (L2)
var ErrTxPoolIDNotFound = fmt.Errorf("%w: transaction id not found in context", ErrTxPool)

// ErrTxPoolNotFound will indicate that transaction with id that found in context not found in pool
// this is a child error (L2)
var ErrTxPoolNotFound = fmt.Errorf("%w: transaction id not found in pool", ErrTxPool)

// ErrTxPoolTrxStillExistsInPool will indicate that transaction with id that found in context still exists in pool
// this is a child error (L2)
var ErrTxPoolTrxStillExistsInPool = fmt.Errorf("%w: transaction still exists in pool", ErrTxPool)
