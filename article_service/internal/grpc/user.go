package grpc

import (
	"context"
	"fmt"

	"github.com/jrdnull/microservices_technical_test/user_service/userpb"

	"github.com/jrdnull/microservices_technical_test/article_service/internal"
)

// UserRepository provides access to users in the gRPC user service.
type UserRepository struct {
	c userpb.UserServiceClient
}

// NewUserRepository returns a new UserRepository backed by c.
func NewUserRepository(c userpb.UserServiceClient) *UserRepository {
	return &UserRepository{c}
}

// Get user by id.
func (r *UserRepository) Get(ctx context.Context, id int64) (*internal.User, error) {
	resp, err := r.c.GetNewsFeedTags(ctx, &userpb.GetNewsFeedTagsRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("user_service: %w", err)
	}
	return &internal.User{
		ID:           id,
		NewsFeedTags: resp.TagIds,
	}, nil
}
