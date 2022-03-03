package internal

import "context"

// UserRepository provides access to users.
type UserRepository interface {
	Get(ctx context.Context, id int64) (*User, error)
}

// UserService provides methods for working with users.
type UserService struct {
	r UserRepository
}

// NewUserService returns a new UserService backed by r.
func NewUserService(r UserRepository) *UserService {
	return &UserService{r}
}

// Get user by id.
func (s *UserService) Get(ctx context.Context, id int64) (*User, error) {
	return s.r.Get(ctx, id)
}

// User represents a user of the article service.
type User struct {
	ID           int64
	NewsFeedTags []int64
}
