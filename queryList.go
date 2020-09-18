package sqldb

import (
	"context"
)

// QueryList retrieves a list of rows from the database
func (g *gateway) QueryList(ctx context.Context, list interface{}, qry string, args ...interface{}) error {
	conn, err := g.connect(ctx)
	if err != nil {
		return err
	}
	rows, err := conn.QueryContext(ctx, qry, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	return parseRows(rows, list)
}
