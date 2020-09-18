package sqldb

import (
	"context"
)

// QueryItem retrieves a row from the database
func (g *gateway) QueryItem(ctx context.Context, item interface{}, qry string, args ...interface{}) error {
	conn, err := g.connect(ctx)
	if err != nil {
		return err
	}
	row := conn.QueryRowContext(ctx, qry, args...)
	return parseRow(row, item)
}
