package service

// Service ---
type Service struct{ repository Repository }

// NewService ---
func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}
