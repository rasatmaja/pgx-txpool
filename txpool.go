package pgxtxpoll

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TxID an identifier for a transaction
type TxID string

// TxContextID an indentifier for a transaction ID in context
type TxContextID string

// ContextTxKey a key for a transaction ID in context
const ContextTxKey TxContextID = "TX_POLL_ID"

// Pool is a struct that wraps a pgx pool
type Pool struct {
	*pgxpool.Pool
	txpool map[TxID]pgx.Tx
}

// New will create a new connection pgx pool
func New() *Pool {
	var config *pgxpool.Config
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}
	return &Pool{
		Pool:   pool,
		txpool: make(map[TxID]pgx.Tx),
	}
}

// BeginTX will begin a prosgres transaction and create an ID
// then it will save those ID and it tx to the pool
// then inject trx id into context and return it
func (p *Pool) BeginTX(ctx context.Context) (context.Context, error) {

	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// generate tx id
	txID := TxID("123")

	// save tx
	p.txpool[txID] = tx

	ctx = context.WithValue(ctx, ContextTxKey, txID)

	return ctx, nil
}

// CommitTX will commit a transaction
func (p *Pool) CommitTX(ctx context.Context) error {
	txID, ok := ctx.Value(ContextTxKey).(TxID)
	if !ok {
		return fmt.Errorf("no transaction id found in context")
	}

	tx, ok := p.txpool[txID]
	if !ok {
		return fmt.Errorf("transaction id not found in pool")
	}

	delete(p.txpool, txID)
	return tx.Commit(ctx)
}

// RollbackTX will rollback a transaction specific to the context
func (p *Pool) RollbackTX(ctx context.Context) error {
	txID, ok := ctx.Value(ContextTxKey).(TxID)
	if !ok {
		return fmt.Errorf("no transaction id found in context")
	}

	tx, ok := p.txpool[txID]
	if !ok {
		return fmt.Errorf("transaction id not found in pool")
	}

	delete(p.txpool, txID)
	return tx.Rollback(ctx)
}

// Exec will execute a query
// if transaction id is found in context
// then use exec from transaction
// otherwise it will use default exec from pgxpool
func (p *Pool) Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error) {
	// if transaction id is found in context
	// then use exec from transaction
	if txID, ok := ctx.Value(ContextTxKey).(TxID); ok {
		if tx, ok := p.txpool[txID]; ok {
			return tx.Exec(ctx, sql, arguments...)
		}
	}

	// default will use func Exec from pgxpool
	return p.Pool.Exec(ctx, sql, arguments...)
}

// Query will execute a query
// if transaction id is found in context
// then use query from transaction
// otherwise it will use default query from pgxpool
func (p *Pool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	// if transaction id is found in context
	// then use query from transaction
	if txID, ok := ctx.Value(ContextTxKey).(TxID); ok {
		if tx, ok := p.txpool[txID]; ok {
			return tx.Query(ctx, sql, args...)
		}
	}

	// default will use func Query from pgxpool
	return p.Pool.Query(ctx, sql, args...)
}

// VerifyTX will verify a transaction to make sure it is not in the pool
// and transaction corelated with this context already commit or rollback
// use this function after using BeginTX
// ex: defer p.VerifyTX(ctx)
func (p *Pool) VerifyTX(ctx context.Context) error {
	if txID, ok := ctx.Value(ContextTxKey).(TxID); ok {
		// if transaction id is found in context
		// return error
		if _, ok := p.txpool[txID]; ok {
			return fmt.Errorf("transaction id still exists in pool")
		}
	}
	return nil
}
