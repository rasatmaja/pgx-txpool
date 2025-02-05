package service

import (
	"context"

	"github.com/rasatmaja/pgx-txpool/tests/integration/model"
)

// TransferBalance --
func (s *Service) TransferBalance(ctx context.Context, trx []model.Transaction, transfer model.TransactionTransfer) error {

	trxCTX, err := s.repository.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func(cause error) {
		if cause != nil {
			err = s.repository.RollbackTx(trxCTX)
		}
	}(err)

	if err = s.repository.CreateTransaction(trxCTX, trx); err != nil {
		return err
	}

	if err = s.repository.CreateTransactionTransfer(trxCTX, []model.TransactionTransfer{transfer}); err != nil {
		return err
	}

	for _, transaction := range trx {
		transferSign := 1.0
		if transaction.Type == "TRANSFER_OUT" {
			transferSign = -1.0
		}

		if err = s.repository.UpdateUserBalance(trxCTX, model.User{
			ID:            transaction.UserID,
			BalanceChange: transaction.Amount * transferSign,
		}); err != nil {
			return err
		}
	}

	if err = s.repository.CommitTx(trxCTX); err != nil {
		return err
	}

	return nil
}
