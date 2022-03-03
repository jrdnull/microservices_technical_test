package main

import (
	"context"
	"log"
	"os"

	"github.com/uptrace/bun/dbfixture"

	"github.com/jrdnull/microservices_technical_test/user_service/cmd"
	"github.com/jrdnull/microservices_technical_test/user_service/internal/postgres"
)

func main() {
	cfg, err := cmd.NewConfig()
	if err != nil {
		log.Fatalf("internal.NewConfig: %v", err)
	}

	db, err := postgres.NewDB(cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port)
	if err != nil {
		log.Fatalf("postgres.NewDB: %v", err)
	}

	fixture := dbfixture.New(db, dbfixture.WithRecreateTables())
	err = fixture.Load(context.Background(), os.DirFS("./cmd/initdb"), "fixture.yaml")
	if err != nil {
		log.Fatalf("fixture.Load: %v", err)
	}
}
