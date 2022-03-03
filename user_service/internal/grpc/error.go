package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jrdnull/microservices_technical_test/user_service/internal"
)

func recoveryHandler(ctx context.Context, r interface{}) error {
	var err error
	switch r := r.(type) {
	case error:
		err = r
	default:
		err = fmt.Errorf("%v", r)
	}

	// TODO: log with context to Sentry etc
	log.Printf("panic recovered: %v", err)

	return status.Error(codes.Internal, err.Error())
}

// statusError returns a suitable status.Error based on the error cause.
func statusError(err error) error {
	if err == nil {
		return nil
	}

	cause, code := err, codes.Internal
	for cause != nil {
		if err, ok := cause.(*internal.Error); ok {
			switch err.Code {
			case internal.ErrorCodeValidation:
				code = codes.InvalidArgument
			}
		} else if cause == context.DeadlineExceeded {
			code = codes.DeadlineExceeded
		}
		cause = errors.Unwrap(cause)
	}
	return status.Error(code, err.Error())
}
