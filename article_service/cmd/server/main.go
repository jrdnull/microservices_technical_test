package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jrdnull/microservices_technical_test/user_service/userpb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"

	"github.com/jrdnull/microservices_technical_test/article_service/cmd"
	"github.com/jrdnull/microservices_technical_test/article_service/internal"
	igrpc "github.com/jrdnull/microservices_technical_test/article_service/internal/grpc"
	"github.com/jrdnull/microservices_technical_test/article_service/internal/http"
	"github.com/jrdnull/microservices_technical_test/article_service/internal/jaeger"
	"github.com/jrdnull/microservices_technical_test/article_service/internal/postgres"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := cmd.NewConfig()
	if err != nil {
		return fmt.Errorf("internal.NewConfig: %w", err)
	}

	tp, err := jaeger.NewTracer(cfg.JaegerEndpoint)
	if err != nil {
		return fmt.Errorf("jaeger.NewTracer: %v", err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("shutdown trace provider: %v", err)
		}
	}()

	db, err := postgres.NewDB(cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port)
	if err != nil {
		return fmt.Errorf("postgres.NewDB: %w", err)
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	}
	conn, err := grpc.Dial(cfg.UserService.Address, opts...)
	if err != nil {
		return fmt.Errorf("grpc.Dial: %w", err)
	}
	defer conn.Close() // nolint: errcheck

	server := http.NewServer(
		internal.NewArticleService(postgres.NewArticleRepository(db)),
		internal.NewUserService(igrpc.NewUserRepository(userpb.NewUserServiceClient(conn))),
		http.Config{
			ArticleTagFilterInputs: cfg.ArticleTagFilterInputs,
		},
	)
	errc := make(chan error, 1)
	go func() {
		errc <- server.Start(fmt.Sprintf(":%d", cfg.Port))
	}()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case <-quit:
		log.Print("shutting down server gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutdown echo: %w", err)
		}
		log.Print("shutdown complete!")
		return nil
	case err := <-errc:
		return err
	}
}
