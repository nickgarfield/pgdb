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
	rows, err := conn.QueryContext(ctx, qry, args...)
	return g.parseRow(rows, item)
}
