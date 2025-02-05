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