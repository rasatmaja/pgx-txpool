package service

import (
	"context"

	"github.com/rasatmaja/pgx-txpool/tests/integration/model"
)

// Repository ---
type Repository interface {
	BeginTx(ctx context.Context) (context.Context, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error
	VerifyTX(ctx context.Context) error
	CreateUser(ctx context.Context, user model.User) error
	GetUsers(ctx context.Context) ([]model.User, error)
	UpdateUserBalance(ctx context.Context, user model.User) error
	CreateTransaction(ctx context.Context, transactions []model.Transaction) error
	CreateTransactionTransfer(ctx context.Context, transactions []model.TransactionTransfer) error
}
