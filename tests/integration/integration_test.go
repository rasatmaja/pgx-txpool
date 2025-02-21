//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	pgxtxpool "github.com/rasatmaja/pgx-txpool"
	"github.com/rasatmaja/pgx-txpool/tests/integration/model"
	"github.com/rasatmaja/pgx-txpool/tests/integration/repository"
	"github.com/rasatmaja/pgx-txpool/tests/integration/service"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestSuite struct {
	repo *repository.Repository
	srv  *service.Service
}

func TestMain(t *testing.T) {
	ctx := context.Background()

	suite := &TestSuite{}
	// setup test containers
	assert.NotPanics(t, func() { suite.Setup(ctx) })

	// run tests
	t.Run("TestMigration", suite.Migration)
	t.Run("TestCreateUser", suite.CreateUser)
}

func (ts *TestSuite) Setup(ctx context.Context) {
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
			WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(10 * time.Second),
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

	ts.repo = repository.NewRepository(db)
	ts.srv = service.NewService(ts.repo)
}

// Migration tests repository Migration method
func (ts *TestSuite) Migration(t *testing.T) {
	ctx := context.Background()
	err := ts.repo.Migration(ctx)
	assert.NoError(t, err, "failed execute migration")

	columnsUsers, err := ts.repo.ShowColomns(ctx, "users")
	assert.NoError(t, err, "failed to get columns users")
	assert.ElementsMatch(t, []string{"id", "name", "balance"}, columnsUsers)

	columnsTransactions, err := ts.repo.ShowColomns(ctx, "transactions")
	assert.NoError(t, err, "failed to get columns transactions")
	assert.ElementsMatch(t, []string{"id", "user_id", "type", "amount"}, columnsTransactions)

	columnsTransactionsTransfer, err := ts.repo.ShowColomns(ctx, "transactions_transfer")
	assert.NoError(t, err, "failed to get columns transactions_transfer")
	assert.ElementsMatch(t, []string{"id", "transaction_origin_id", "transaction_destination_id", "amount"}, columnsTransactionsTransfer)
}

// CreateUser tests service CreateUser method
func (ts *TestSuite) CreateUser(t *testing.T) {

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

			err := ts.srv.CreateUser(ctx, c.user, c.trx...)
			// TODO: Should assert specific error
			assert.Equal(t, c.error, err != nil)
		})
	}

	t.Run("check data integrity", func(t *testing.T) {
		t.Skip("not implemented yet")
	})

}
