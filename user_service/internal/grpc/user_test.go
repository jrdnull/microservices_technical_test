package grpc

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jrdnull/microservices_technical_test/user_service/internal"
	"github.com/jrdnull/microservices_technical_test/user_service/internal/mock"
	"github.com/jrdnull/microservices_technical_test/user_service/userpb"
)

func TestUserServiceServerGetNewsFeedTags(t *testing.T) {
	cases := map[string]struct {
		req               *userpb.GetNewsFeedTagsRequest
		getNewsFeedTagsFn func(context.Context, int64) ([]int64, error)
		wantResp          *userpb.GetNewsFeedTagsResponse
		wantErr           error
	}{
		"missing id": {
			req:               &userpb.GetNewsFeedTagsRequest{},
			getNewsFeedTagsFn: nil,
			wantResp:          nil,
			wantErr:           status.Error(codes.InvalidArgument, "missing required field: Id"),
		},
		"panic": {
			req: &userpb.GetNewsFeedTagsRequest{Id: 1},
			getNewsFeedTagsFn: func(ctx context.Context, id int64) ([]int64, error) {
				panic("oops")
			},
			wantResp: nil,
			wantErr:  status.Error(codes.Internal, "oops"),
		},
		"repo error": {
			req: &userpb.GetNewsFeedTagsRequest{Id: 1},
			getNewsFeedTagsFn: func(ctx context.Context, id int64) ([]int64, error) {
				return nil, errors.New("something went wrong")
			},
			wantResp: nil,
			wantErr:  status.Error(codes.Internal, "something went wrong"),
		},
		"success": {
			req: &userpb.GetNewsFeedTagsRequest{Id: 1},
			getNewsFeedTagsFn: func(ctx context.Context, id int64) ([]int64, error) {
				return []int64{1, 2, 3}, nil
			},
			wantResp: &userpb.GetNewsFeedTagsResponse{TagIds: []int64{1, 2, 3}},
			wantErr:  nil,
		},
	}
	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			client, cleanup := testServer(t, internal.NewUserService(mock.UserRepository{
				GetNewsFeedTagsFn: tc.getNewsFeedTagsFn,
			}))
			defer cleanup()

			ctx := context.Background()
			resp, err := client.GetNewsFeedTags(ctx, tc.req)
			if tc.wantErr != nil {
				if diff := cmp.Diff(err, tc.wantErr, cmpopts.EquateErrors()); diff != "" {
					t.Fatalf("(-got +want):\n%s", diff)
				}
				return // test passed
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			diff := cmp.Diff(resp, tc.wantResp, cmpopts.IgnoreUnexported(userpb.GetNewsFeedTagsResponse{}))
			if diff != "" {
				t.Errorf("(-got +want):\n%s", diff)
			}
		})
	}
}

func testServer(t *testing.T, users *internal.UserService) (client userpb.UserServiceClient, cleanup func()) {
	t.Helper()

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("net.Listen: %v", err)
	}

	server := NewServer(users)
	go server.Serve(l) // nolint: errcheck

	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(l.Addr().String(), opts...)
	if err != nil {
		t.Fatalf("grpc.Dial: %v", err)
	}

	return userpb.NewUserServiceClient(conn), func() {
		_ = conn.Close()
		server.Stop()
		_ = l.Close()
	}
}
