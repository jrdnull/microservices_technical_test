package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/jrdnull/microservices_technical_test/user_service/cmd"
	"github.com/jrdnull/microservices_technical_test/user_service/internal"
	"github.com/jrdnull/microservices_technical_test/user_service/internal/grpc"
	"github.com/jrdnull/microservices_technical_test/user_service/internal/jaeger"
	"github.com/jrdnull/microservices_technical_test/user_service/internal/postgres"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := cmd.NewConfig()
	if err != nil {
		return fmt.Errorf("internal.NewConfig: %v", err)
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
		return fmt.Errorf("postgres.NewDB: %v", err)
	}

	server := grpc.NewServer(internal.NewUserService(postgres.NewUserRepository(db)))
	errc := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Port)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			errc <- fmt.Errorf("net.Listen: %w", err)
			return
		}
		log.Printf("starting gRPC server on %s", addr)
		errc <- server.Serve(lis)
	}()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case <-quit:
		log.Print("shutting down server gracefully...")
		server.GracefulStop()
		log.Print("shutdown complete!")
		return nil
	case err := <-errc:
		return err
	}
}
