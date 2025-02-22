package service

import (
	"context"

	"github.com/rasatmaja/pgx-txpool/tests/integration/model"
)

// ListUser --
func (s *Service) ListUser(ctx context.Context) ([]model.User, error) {
	return s.repository.GetUsers(ctx)
}
