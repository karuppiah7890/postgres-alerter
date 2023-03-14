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

type Client struct {
	db *bun.DB
}

func NewClient(uri string) *Client {
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))
	db := bun.NewDB(pgdb, pgdialect.New())
	return &Client{
		db: db,
	}
}

func (c *Client) GetPostgresStatus() PostgresStatus {
	ctx := context.Background()

	err := c.db.PingContext(ctx)
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
