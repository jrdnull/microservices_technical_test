FROM golang:1.16-alpine AS build
WORKDIR /go/src/article_service
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/server ./cmd/server

FROM scratch
COPY --from=build /go/bin/server /bin/server
ENTRYPOINT ["/bin/server"]
