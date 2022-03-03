package grpc

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"github.com/jrdnull/microservices_technical_test/user_service/internal"
	"github.com/jrdnull/microservices_technical_test/user_service/userpb"
)

// NewServer returns a new grpc.Server with services registered.
func NewServer(users *internal.UserService) *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandlerContext(recoveryHandler),
			),
			otelgrpc.UnaryServerInterceptor(),
		)),
	)

	userpb.RegisterUserServiceServer(s, &userServiceServer{users: users})
	return s
}
