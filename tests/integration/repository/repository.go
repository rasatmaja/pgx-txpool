package repository

import (
	"context"

	pgxtxpool "github.com/rasatmaja/pgx-txpool"
	"github.com/rasatmaja/pgx-txpool/tests/integration/model"
)

// Repository --
type Repository struct {
	db *pgxtxpool.Pool
}

// NewRepository ---
func NewRepository(db *pgxtxpool.Pool) *Repository {
	return &Repository{db: db}
}

// BeginTx ---
func (r *Repository) BeginTx(ctx context.Context) (context.Context, error) {
	return r.db.BeginTX(ctx)
}

// CommitTx ---
func (r *Repository) CommitTx(ctx context.Context) error {
	return r.db.CommitTX(ctx)
}

// RollbackTx ---
func (r *Repository) RollbackTx(ctx context.Context) error {
	return r.db.RollbackTX(ctx)
}

// VerifyTX --
func (r *Repository) VerifyTX(ctx context.Context) error {
	return r.db.VerifyTX(ctx)
}

// CreateUser ---
func (r *Repository) CreateUser(ctx context.Context, user model.User) error {
	query := `INSERT INTO users (id, name, balance) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Balance)
	if err != nil {
		return err
	}
	return nil
}

// GetUsers --
func (r *Repository) GetUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	query := `SELECT id, name, balance FROM users`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Balance)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// UpdateUserBalance ---
func (r *Repository) UpdateUserBalance(ctx context.Context, user model.User) error {
	query := `UPDATE users SET balance = balance + $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, user.BalanceChange, user.ID)
	if err != nil {
		return err
	}
	return nil
}

// CreateTransaction ---
func (r *Repository) CreateTransaction(ctx context.Context, transactions []model.Transaction) error {
	query := `INSERT INTO transactions (id, user_id, type, amount) VALUES ($1, $2, $3, $4)`
	for _, transaction := range transactions {
		_, err := r.db.Exec(ctx, query, transaction.ID, transaction.UserID, transaction.Type, transaction.Amount)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateTransactionTransfer ---
func (r *Repository) CreateTransactionTransfer(ctx context.Context, transactions []model.TransactionTransfer) error {
	query := `INSERT INTO transaction_transfers (id, transaction_origin_id, transaction_destination_id ,amount) VALUES ($1, $2, $3)`
	for _, transaction := range transactions {
		_, err := r.db.Exec(ctx, query, transaction.ID, transaction.TransactionOriginID, transaction.TransactionDestinationID, transaction.Amount)
		if err != nil {
			return err
		}
	}
	return nil
}
