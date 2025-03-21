# pgx-txpool
A Go package that provides transaction pooling functionality for PostgreSQL using the pgx driver. This package provides a convenient way to manage database transactions with context-based transaction tracking.

## Background
Why this package exists? This package is created to solve several common challenges in database transaction management:

### DDD Repository Pattern Implementation
When implementing Domain-Driven Design (DDD), it's common to have one function per table/context in repositories. This often leads to creating two versions of the same function:
- One that accepts a transaction parameter for multi-operation transactions.
```go
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, data *Data) error {
    // ...
}
```

- Another that uses a regular database connection for single operations. 
```go
func (r *Repository) Create(ctx context.Context, data *Data) error {
    //...
}
```
This duplication increases code maintenance overhead and complexity.

> ***why not write a single function that accepts only transaction connection?*** 

I think that not efficient, because when we need only one operation, For example, we just performing a simple insert, the overhead of creating and managing a transaction might be unnecessary.

> ***how about using pattern like Repository Pattern with Query Object ?***
```go
// ISQL is a interface for database/sql
type ISQL interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

//ISQLTX is a interface for sql transaction
type ISQLTX interface {
	Rollback() error
	Commit() error
}

// Queries ...
type Queries struct {
	DB ISQL
}

func New() *Queries {
    // open connection to PostgreSQL server
	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		return nil, fmt.Errorf("cannot open connection, got: %v", err)
	}
	return &Queries{
		DB: db,
	}
}

// BeginTX is a function to start transaction
func (qry *Queries) BeginTX(ctx context.Context) (*Queries, ISQLTX, error) {
	db, err := qry.OpenConn()
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to open connection, got: %v", err)
	}

	// begin transactions
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to open transaction, got: %v", err)
	}
	queries := &Queries{DB: tx}

	return queries, tx, nil
}

// The DB in this case would be either a regular database connection or a transaction connection, it depends on whether a queries is created with a transaction (BeginTX) or a regular connection (New).
func (qry *Queries) Create(ctx context.Context, data *Data) error {
    qry.DB.ExecContext(ctx, "INSERT INTO table_name (column1, column2) VALUES ($1, $2)", data.Column1, data.Column2)
    //...
}
```
I think this is a good solution, but it's not perfect. Because we need to create a new struct (Queries) for each transaction operation. The other problem is, if we have multiple domain and multiple repository, we need to create a new struct for each domain and repository. 

Sometime in one domain we need to use repository from another domain and repository from another domain is used for part of trasaction, for example:
```go
// users/repository.go -> User Domain
type UserRepository struct {
    DB ISQL
}

// transactions/repository.go -> Transaction Domain
type TransactionRepository struct {
    DB ISQL
}

// trasactions/service.go -> Transaction Domain
type TransactionService struct {
    UserRepo *UserRepository
    TransactionRepo *TransactionRepository
}

func (svc *TransactionService) CreateTransaction(ctx context.Context, data *Transaction) error {
    //...
    // we need to use UserRepo to update user data
    repoUsr, txUsr, err := svc.UserRepo.BeginTX(ctx)
    user, err := repoUsr.UpdateBalance(ctx, data)
    if err!= nil {
        txUsr.Rollback()
        return err
    }
    // we need to use TransactionRepo to create transaction data
    repoTx, txTx, err := svc.TransactionRepo.BeginTX(ctx)
    transaction, err := repoTx.Create(ctx, data)
    if err!= nil {
        txTx.Rollback()
        return err
    }

    txUsr.Commit()
    txTx.Commit()
    //...
}

```

From this example, we can see that we need to create a new struct for each domain and repository. And i think that it's not good solution, because epproach has two main issues:

1. **Multiple Transaction Management**
When dealing with multiple domains (like User and Transaction), each repository creates its own separate transaction. This means:
    - Each transaction has a different ID
    - Transactions are not coordinated with each other
    - We need to manage multiple commit/rollback operations

2. **Risk of Transaction Leaks**
Consider this scenario:
```go
// First transaction for updating user balance
repoUsr.UpdateBalance()  ✅ Success

// Second transaction for creating transaction record
repoTx.Create()  ❌ Failed
```
If the second operation fails, we must remember to:
- Rollback the transaction record
- Rollback the user balance update
Forgetting to rollback any of these transactions will leave them hanging in the database, causing resource leaks and potential data inconsistency.

This complexity grows exponentially as we add more domains and repositories, making it error-prone and difficult to maintain.

### Context-Based Transaction Management
Some projects use context to manage transactions by storing the transaction connection in the context object. Here's how it typically works:
```go

// users/repository.go -> User Domain
type UserRepository struct {}

func (repo *UserRepository) UpdateBalance(ctx context.Context, data *Data) error {
    //...
    // we need to check if transaction connection is exist in context
    tx, ok := ctx.Value("tx").(pgx.Tx)
    if ok {
        tx.ExecContext(ctx, "UPDATE users SET balance = $1 WHERE id = $2", data.Balance, data.ID)
    } else {
        repo.DB.ExecContext(ctx, "UPDATE users SET balance = $1 WHERE id = $2", data.Balance, data.ID)
    }
    //...
}
```

This approach has several drawbacks:

1. Context Pollution  
   - The context object becomes a carrier for database connections
   - This goes against the intended use of context for request cancellation and deadlines
   - Makes it harder to track what's being stored in the context
2. Type Safety Issues
   - Requires type assertion ( `ctx.Value("tx").(pgx.Tx)` )
   - Can lead to runtime errors if the wrong type is stored
   - Makes the code more fragile and harder to maintain
3. Hidden Dependencies
   - The transaction dependency is not explicit in the function signature
   - Makes it harder to understand what the function needs to work properly

## Features
- Transaction pooling with context-based tracking
- Transaction pooling management
- Safe concurrent transaction handling
- Transaction verification to prevent leaks

## Requirements
- Go 1.21 or higher
- PostgreSQL
- jackc/pgx/v5 driver

## Installation

```bash
go get github.com/rasatmaja/pgx-txpool
```

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.