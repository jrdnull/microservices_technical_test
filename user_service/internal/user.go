package internal

import (
	"context"
	"net/mail"
)

// UserRepository provides access to users.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	SetNewsFeedTags(ctx context.Context, id int64, tags []int64) error
	GetNewsFeedTags(ctx context.Context, id int64) ([]int64, error)
}

// UserService provides methods for working with users.
type UserService struct {
	r UserRepository
}

// NewUserService returns a new UserService backed by r.
func NewUserService(r UserRepository) *UserService {
	return &UserService{r}
}

// Create user.
// (not actually used just an example of business logic is contained in
// the service layer)
func (s *UserService) Create(ctx context.Context, user *User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	return s.r.Create(ctx, user)
}

// SetNewsFeedTags sets the tag ids for the users news feed.
func (s *UserService) SetNewsFeedTags(ctx context.Context, id int64, tagIDs []int64) error {
	return s.r.SetNewsFeedTags(ctx, id, tagIDs)
}

// GetNewsFeedTags returns the users tag ids for their news feed.
func (s *UserService) GetNewsFeedTags(ctx context.Context, id int64) ([]int64, error) {
	return s.r.GetNewsFeedTags(ctx, id)
}

// User represents a user of the news feed system.
type User struct {
	ID    int64  `bun:",pk,autoincrement"`
	Email string `bun:",unique,notnull"`
}

// Validate the receiver.
func (u *User) Validate() error {
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return NewValidationError("invalid email")
	}
	return nil
}
