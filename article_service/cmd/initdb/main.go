package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"text/template"
	"time"

	"github.com/uptrace/bun/dbfixture"

	"github.com/jrdnull/microservices_technical_test/article_service/cmd"
	"github.com/jrdnull/microservices_technical_test/article_service/internal/postgres"
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

	funcMap := template.FuncMap{
		"randomTimestamp": func() string {
			return time.Now().Add(time.Duration(-rand.Intn(43800)) * time.Minute).Format(time.RFC3339Nano)
		},
	}
	fixture := dbfixture.New(db, dbfixture.WithRecreateTables(), dbfixture.WithTemplateFuncs(funcMap))
	err = fixture.Load(context.Background(), os.DirFS("./cmd/initdb"), "fixture.yaml")
	if err != nil {
		log.Fatalf("fixture.Load: %v", err)
	}
}
