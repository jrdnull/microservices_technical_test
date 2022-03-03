package grpc

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jrdnull/microservices_technical_test/user_service/internal"
)

func TestStatusError(t *testing.T) {
	tests := map[string]struct {
		in   error
		want *status.Status
	}{
		"nil":           {nil, nil},
		"regular error": {errors.New("test"), status.New(codes.Internal, "test")},
		"internal.Error": {
			internal.NewValidationError("test"),
			status.New(codes.InvalidArgument, "test"),
		},
		"wrapped internal.Error": {
			fmt.Errorf("b: %w", internal.NewValidationError("a")),
			status.New(codes.InvalidArgument, "b: a"),
		},
		"context.DeadlineExceeded": {
			context.DeadlineExceeded,
			status.New(codes.DeadlineExceeded, context.DeadlineExceeded.Error()),
		},
	}

	for desc, tt := range tests {
		t.Run(desc, func(t *testing.T) {
			got := statusError(tt.in)
			if tt.want == nil {
				if got != nil {
					t.Fatalf("unexpected error: %v", got)
				}
				return
			}

			statusErr, ok := status.FromError(got)
			if !ok {
				t.Fatalf("expected status error, got: %T", got)
			}

			if statusErr.Code() != tt.want.Code() {
				t.Errorf("code, got: %v, want: %v", statusErr.Code(), tt.want.Code())
			}
			if statusErr.Message() != tt.want.Message() {
				t.Errorf("message, got: %q, want: %q", statusErr.Message(), tt.want.Message())
			}
		})
	}
}
