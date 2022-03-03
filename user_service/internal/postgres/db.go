package postgres

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bunotel"

	"github.com/jrdnull/microservices_technical_test/user_service/internal"
)

// NewDB opens, verifies and returns new bun.DB.
func NewDB(user, pass, host string, port int) (*bun.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/user_service?sslmode=disable",
		user, pass, host, port)
	db := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn))), pgdialect.New())
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	db.AddQueryHook(bunotel.NewQueryHook())

	tables := []interface{}{
		(*internal.User)(nil),
		(*UserNewsFeed)(nil),
	}
	for _, v := range tables {
		db.RegisterModel(v)
	}

	return db, nil
}
