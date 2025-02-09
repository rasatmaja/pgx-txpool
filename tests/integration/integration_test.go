package integration

import (
	"context"
	"os"
	"testing"

	pgxtxpool "github.com/rasatmaja/pgx-txpool"
	"github.com/rasatmaja/pgx-txpool/tests/integration/model"
	"github.com/rasatmaja/pgx-txpool/tests/integration/repository"
	"github.com/rasatmaja/pgx-txpool/tests/integration/service"
)

var repo *repository.Repository
var srv *service.Service

func TestMain(m *testing.M) {
	// setup database
	db := pgxtxpool.New()
	repo = repository.NewRepository(db)
	srv = service.NewService(repo)
	os.Exit(m.Run())
}

// TestCreateUser tests service CreateUser method
func TestCreateUser(t *testing.T) {
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
			if err != nil {
				t.Errorf("failed to create user: %v", err)
			}
		})
	}

	t.Run("check data integrity", func(t *testing.T) {
		t.Parallel()
	})

}
