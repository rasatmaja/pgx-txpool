//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"testing"

	pgxtxpool "github.com/rasatmaja/pgx-txpool"
	"github.com/rasatmaja/pgx-txpool/tests/integration/model"
	"github.com/rasatmaja/pgx-txpool/tests/integration/repository"
	"github.com/rasatmaja/pgx-txpool/tests/integration/service"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var repo *repository.Repository
var srv *service.Service

func TestMain(t *testing.T) {
	ctx := context.Background()

	// setup test containers
	SetupTest(ctx)

	// run tests
	t.Run("TestMigration", TestMigration)
	t.Run("TestCreateUser", TestCreateUser)
}

func SetupTest(ctx context.Context) {
	var pgUsername, pgPassword, pgDatabase, pgHost, pgPort string
	pgUsername = "postgres-user"
	pgPassword = "postgres-password"
	pgDatabase = "postgres-db"

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{"5432/tcp"},
			AutoRemove:   true,
			Env: map[string]string{
				"POSTGRES_USER":     pgUsername,
				"POSTGRES_PASSWORD": pgPassword,
				"POSTGRES_DB":       pgDatabase,
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	if err != nil {
		panic(err)
	}

	// get hostname from generated test container
	pgHost, err = postgres.Host(ctx)
	if err != nil {
		panic(err)
	}

	// get exposed port from generated test container
	exposedPort, err := postgres.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}

	pgPort = exposedPort.Port()

	// setup database
	db := pgxtxpool.New(
		pgxtxpool.SetHost(pgHost, pgPort),
		pgxtxpool.SetCredential(pgUsername, pgPassword),
		pgxtxpool.SetDatabase(pgDatabase),
		pgxtxpool.WithSSLMode("disable"),
		pgxtxpool.WithMaxConns(20),
		pgxtxpool.WithMaxIdleConns("30s"),
		pgxtxpool.WithMaxConnLifetime("5m"),
	)

	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	repo = repository.NewRepository(db)
	srv = service.NewService(repo)
}

// TestMigration tests repository Migration method
func TestMigration(t *testing.T) {
	ctx := context.Background()
	if err := repo.Migration(ctx); err != nil {
		t.Errorf("failed to migrate: %v", err)
		t.FailNow()
	}

	columnsUsers, err := repo.ShowColomns(ctx, "users")
	if err != nil {
		t.Errorf("failed to get columns: %v", err)
	}

	t.Log(fmt.Sprintf("table users: %v", columnsUsers))
	if len(columnsUsers) != 3 {
		t.Errorf("expected 3 columns, got %d", len(columnsUsers))
	}

	columnsTransactions, err := repo.ShowColomns(ctx, "transactions")
	if err != nil {
		t.Errorf("failed to get columns: %v", err)
	}

	t.Log(fmt.Sprintf("table transactions: %v", columnsTransactions))
	if len(columnsTransactions) != 4 {
		t.Errorf("expected 4 columns, got %d", len(columnsTransactions))
	}

	columnsTransactionsTransfer, err := repo.ShowColomns(ctx, "transactions_transfer")
	if err != nil {
		t.Errorf("failed to get columns: %v", err)
	}

	t.Log(fmt.Sprintf("table transactions_transfer: %v", columnsTransactionsTransfer))
	if len(columnsTransactionsTransfer) != 4 {
		t.Errorf("expected 4 columns, got %d", len(columnsTransactionsTransfer))
	}
}

// TestCreateUser tests service CreateUser method
func TestCreateUser(t *testing.T) {

	if repo == nil || srv == nil {
		t.Skip("test skipped cause repository or service is nil")
	}

	ctx := context.Background()
	cases := []struct {
		name  string
		error bool
		user  model.User
		trx   []model.Transaction
	}{
		{
			name:  "USR001: should success",
			error: false,
			user: model.User{
				ID:      "USR001",
				Name:    "John Doe",
				Balance: 1000,
			},
			trx: []model.Transaction{
				{
					ID:     "TRX001",
					UserID: "USR001",
					Type:   "INITIAL_BALANCE",
					Amount: 1000,
				},
			},
		},
		{
			name:  "USR002: should rollback user data when create trx error",
			error: true,
			user: model.User{
				ID:      "USR002",
				Name:    "Jane",
				Balance: 2000,
			},
			trx: []model.Transaction{
				{
					ID:     "TRX002",
					UserID: "XXXXX",
					Type:   "INITIAL_BALANCE",
					Amount: 2000,
				},
			},
		},
		{
			name:  "USR003: should success",
			error: false,
			user: model.User{
				ID:      "USR003",
				Name:    "Peterson",
				Balance: 3000,
			},
			trx: []model.Transaction{
				{
					ID:     "TRX003",
					UserID: "USR003",
					Type:   "INITIAL_BALANCE",
					Amount: 3000,
				},
			},
		},
		{
			name:  "USR004: should success",
			error: false,
			user: model.User{
				ID:      "USR004",
				Name:    "Waller",
				Balance: 4000,
			},
			trx: []model.Transaction{
				{
					ID:     "TRX004",
					UserID: "USR004",
					Type:   "INITIAL_BALANCE",
					Amount: 4000,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			err := srv.CreateUser(ctx, c.user, c.trx...)
			// TODO: Should assert specific error
			if err != nil && !c.error {
				t.Errorf("failed to create user: %v", err)
			}
		})
	}

	t.Run("check data integrity", func(t *testing.T) {
		t.Skip("not implemented yet")
	})

}
