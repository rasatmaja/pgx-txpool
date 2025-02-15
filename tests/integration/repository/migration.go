package repository

import (
	"context"
	"strings"
)

// Migration ---
func (r *Repository) Migration(ctx context.Context) error {
	query := `
	DROP TABLE IF EXISTS users;
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT,
		balance NUMERIC(10, 2)
	);

	DROP TABLE IF EXISTS transactions;
	CREATE TABLE IF NOT EXISTS transactions (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		type TEXT,
		amount NUMERIC(10, 2),
		CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
	);

	DROP TABLE IF EXISTS transactions_transfer;
	CREATE TABLE IF NOT EXISTS transactions_transfer (
		id TEXT PRIMARY KEY,
		transaction_origin_id TEXT,
		transaction_destination_id TEXT,
		amount NUMERIC(10, 2),
		CONSTRAINT fk_transaction_origin FOREIGN KEY (transaction_origin_id) REFERENCES transactions(id),
		CONSTRAINT fk_transaction_destination FOREIGN KEY (transaction_destination_id) REFERENCES transactions(id)
	);
	`
	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

// ShowColomns --
func (r *Repository) ShowColomns(ctx context.Context, table string) ([]string, error) {
	var rawColumns string
	var columns []string
	query := `SELECT STRING_AGG(column_name, ';') FROM information_schema.columns WHERE table_schema = 'public' AND table_name = $1;`
	row := r.db.QueryRow(ctx, query, table)
	err := row.Scan(&rawColumns)
	if err != nil {
		return columns, err
	}
	columns = strings.Split(rawColumns, ";")
	return columns, nil
}
