package pgdb

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

// Gateway provides an interface to the PostgreSQL database
type Gateway interface {
	BeginTxn(ctx context.Context) (Txn, error)
	QueryItem(ctx context.Context, item interface{}, qry string, args ...interface{}) error
	QueryList(ctx context.Context, list interface{}, qry string, args ...interface{}) error
}

type gateway struct {
	db *sql.DB
}

// New creates a new database gateway
func New(ctx context.Context, connInfo ConnInfo) (Gateway, error) {
	db, err := sql.Open("postgres", connInfo.String())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to database")
	}
	return &gateway{
		db: db,
	}, err
}
