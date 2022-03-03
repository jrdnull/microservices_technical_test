package mock

import (
	"context"

	"github.com/jrdnull/microservices_technical_test/user_service/internal"
)

// UserRepository is a test internal.UserRepository.
type UserRepository struct {
	internal.UserRepository

	GetNewsFeedTagsFn func(context.Context, int64) ([]int64, error)
}

// GetNewsFeedTags calls GetNewsFeedTagsFn.
func (s UserRepository) GetNewsFeedTags(ctx context.Context, id int64) ([]int64, error) {
	return s.GetNewsFeedTagsFn(ctx, id)
}
