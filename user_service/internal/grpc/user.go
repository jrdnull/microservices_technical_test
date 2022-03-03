package grpc

import (
	"context"

	"github.com/jrdnull/microservices_technical_test/user_service/internal"
	"github.com/jrdnull/microservices_technical_test/user_service/userpb"
)

type userServiceServer struct {
	userpb.UnimplementedUserServiceServer
	users *internal.UserService
}

func (s *userServiceServer) GetNewsFeedTags(ctx context.Context, req *userpb.GetNewsFeedTagsRequest) (*userpb.GetNewsFeedTagsResponse, error) {
	if req.Id < 1 {
		return nil, statusError(internal.NewValidationError("missing required field: Id"))
	}
	tagIDs, err := s.users.GetNewsFeedTags(ctx, req.Id)
	if err != nil {
		return nil, statusError(err)
	}
	return &userpb.GetNewsFeedTagsResponse{TagIds: tagIDs}, nil
}
