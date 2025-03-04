package pgxtxpool

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool is a struct that wraps a pgx pool
type Pool struct {
	*pgxpool.Pool
	txpool     sync.Map
	generateID func() TxID
}

// New will create a new connection pgx pool
func New(opts ...Option) *Pool {
	var config config
	for _, opt := range opts {
		opt(&config)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config.ParseToPGXConfig())
	if err != nil {
		panic(err)
	}
	return &Pool{
		Pool:       pool,
		generateID: generateID,
	}
}

// storeTXConn will store a transaction to the pool
func (p *Pool) storeTXConn(txID TxID, tx pgx.Tx) {
	p.txpool.Store(txID, tx)
}

// getTXConn will get a transaction from the pool (sync.Map)
// then return the transaction (pgx.Tx)
func (p *Pool) getTXConn(txID TxID) (pgx.Tx, bool) {
	tx, ok := p.txpool.Load(txID)
	if ok {
		return tx.(pgx.Tx), ok
	}
	return nil, false
}

func (p *Pool) deleteTXConn(txID TxID) {
	p.txpool.Delete(txID)
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
	txID := p.generateID()

	// save tx
	p.storeTXConn(txID, tx)

	ctx = context.WithValue(ctx, ContextTxKey, txID)

	return ctx, nil
}

// CommitTX will commit a transaction
func (p *Pool) CommitTX(ctx context.Context) error {
	txID, ok := ctx.Value(ContextTxKey).(TxID)
	if !ok {
		return ErrTxPoolIDNotFound
	}

	tx, ok := p.getTXConn(txID)
	if !ok {
		return ErrTxPoolNotFound
	}
	p.deleteTXConn(txID)
	return tx.Commit(ctx)
}

// RollbackTX will rollback a transaction specific to the context
func (p *Pool) RollbackTX(ctx context.Context) error {
	txID, ok := ctx.Value(ContextTxKey).(TxID)
	if !ok {
		return ErrTxPoolIDNotFound
	}

	tx, ok := p.getTXConn(txID)
	if !ok {
		return ErrTxPoolNotFound
	}
	p.deleteTXConn(txID)
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
		if tx, ok := p.getTXConn(txID); ok {
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
		if tx, ok := p.getTXConn(txID); ok {
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
		if _, ok := p.getTXConn(txID); ok {
			return ErrTxPoolTrxStillExistsInPool
		}
	}
	return nil
}
