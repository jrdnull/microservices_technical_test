package mock

import (
	"context"

	"github.com/jrdnull/microservices_technical_test/article_service/internal"
)

// UserRepository is a test internal.UserRepository.
type UserRepository struct {
	internal.UserRepository

	GetFn func(context.Context, int64) (*internal.User, error)
}

// Get calls GetFn.
func (s UserRepository) Get(ctx context.Context, id int64) (*internal.User, error) {
	return s.GetFn(ctx, id)
}
