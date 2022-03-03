module github.com/jrdnull/microservices_technical_test/article_service

go 1.16

replace github.com/jrdnull/microservices_technical_test/user_service => ../user_service

require (
	github.com/DATA-DOG/go-txdb v0.1.4
	github.com/beme/abide v0.0.0-20190723115211-635a09831760
	github.com/caarlos0/env/v6 v6.6.2
	github.com/google/go-cmp v0.5.6
	github.com/kr/text v0.2.0 // indirect
	github.com/labstack/echo/v4 v4.5.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/jrdnull/microservices_technical_test/user_service v0.0.0-00010101000000-000000000000
	github.com/uptrace/bun v0.4.3
	github.com/uptrace/bun/dbfixture v0.4.3
	github.com/uptrace/bun/dialect/pgdialect v0.4.3
	github.com/uptrace/bun/driver/pgdriver v0.4.3
	github.com/uptrace/bun/extra/bunotel v0.4.3
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.22.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.22.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC2
	go.opentelemetry.io/otel/sdk v1.0.0-RC2
	google.golang.org/grpc v1.40.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)
