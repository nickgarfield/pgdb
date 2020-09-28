package sqldb

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// QueryList retrieves a list of rows from the database
func (g *gateway) QueryList(ctx context.Context, list interface{}, qry string, args ...interface{}) error {

	// Validate 'list' is a pointer to a slice
	_, err := isSlicePtr(list)
	if err != nil {
		return err
	}

	// Open connection with tenant isolation
	conn, err := g.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Execute the query
	rows, err := conn.QueryContext(ctx, qry, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Scan rows and return
	return sqlx.StructScan(rows, list)
}
