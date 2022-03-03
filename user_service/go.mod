module github.com/jrdnull/microservices_technical_test/user_service

go 1.16

require (
	github.com/DATA-DOG/go-txdb v0.1.4
	github.com/caarlos0/env/v6 v6.6.2
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/google/go-cmp v0.5.6
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/lib/pq v1.10.2 // indirect
	github.com/uptrace/bun v0.4.3
	github.com/uptrace/bun/dbfixture v0.4.3
	github.com/uptrace/bun/dialect/pgdialect v0.4.3
	github.com/uptrace/bun/driver/pgdriver v0.4.3
	github.com/uptrace/bun/extra/bunotel v0.4.3
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.22.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC2
	go.opentelemetry.io/otel/sdk v1.0.0-RC2
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
)
