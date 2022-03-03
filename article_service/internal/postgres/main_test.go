//go:build integration
// +build integration

package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/caarlos0/env/v6"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	_ "github.com/uptrace/bun/driver/pgdriver"
)

type config struct {
	User string `env:"DB_USER" envDefault:"postgres"`
	Pass string `env:"DB_PASS" envDefault:"dev"`
	Host string `env:"DB_HOST" envDefault:"localhost"`
	Port int    `env:"DB_PORT" envDefault:"5432"`
}

func TestMain(m *testing.M) {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("parse config: %v", err)
	}
	if cfg.Host == "" || cfg.Port == 0 {
		log.Fatal("missing env vars for db, set DB_(HOST|PORT|USER|PASS)")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/article_service?sslmode=disable",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port)
	txdb.Register("txdb", "pg", dsn)

	os.Exit(m.Run())
}

func txDB(t *testing.T) *bun.DB {
	sqldb, err := sql.Open("txdb", fmt.Sprintf("itest_%d", rand.Int()))
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	registerModels(db)
	return db
}
