package sqldb

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// QueryItem retrieves a row from the database
func (g *gateway) QueryItem(ctx context.Context, item interface{}, qry string, args ...interface{}) error {

	// Validate 'item' is a pointer to a struct
	_, err := isStructPtr(item)
	if err != nil {
		return err
	}

	// Open connection with tenant isolation
	conn, err := g.connect(ctx)
	if err != nil {
		return err
	}

	// Query the rows
	rows, err := conn.QueryContext(ctx, qry, args...)
	if err != nil {
		return err
	}

	// Scan row and return
	rows.Next()
	return sqlx.StructScan(rows, item)
}
