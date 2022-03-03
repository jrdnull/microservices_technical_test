package postgres

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bunotel"

	"github.com/jrdnull/microservices_technical_test/article_service/internal"
)

// NewDB opens, verifies and returns new bun.DB.
func NewDB(user, pass, host string, port int) (*bun.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/article_service?sslmode=disable",
		user, pass, host, port)
	db := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn))), pgdialect.New())
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	db.AddQueryHook(bunotel.NewQueryHook())

	registerModels(db)

	return db, nil
}

func registerModels(db *bun.DB) {
	tables := []interface{}{
		(*ArticleTag)(nil),
		(*internal.Article)(nil),
		(*internal.Tag)(nil),
	}
	for _, v := range tables {
		db.RegisterModel(v)
	}
}
