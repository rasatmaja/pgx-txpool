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
	user := model.User{
		ID:      "123",
		Name:    "test",
		Balance: 100,
	}
	// TODO: test case user goes here
	t.Run("should XOXO", func(t *testing.T) {
		t.Parallel()

		err := srv.CreateUser(ctx, user)
		if err != nil {
			t.Errorf("failed to create user: %v", err)
		}
	})

	t.Run("check data integrity", func(t *testing.T) {
		t.Parallel()
	})

}
