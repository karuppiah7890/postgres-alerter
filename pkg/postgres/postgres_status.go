package postgres

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type PostgresStatus struct {
	IsUp   bool
	Errors []error
}

func GetPostgresStatus(uri string) PostgresStatus {
	ctx := context.Background()
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))
	db := bun.NewDB(pgdb, pgdialect.New())

	err := db.PingContext(ctx)
	if err != nil {
		return PostgresStatus{
			IsUp:   false,
			Errors: []error{err},
		}
	}

	return PostgresStatus{
		IsUp: true,
	}
}
