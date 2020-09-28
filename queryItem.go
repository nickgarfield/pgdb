package sqldb

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/jmoiron/sqlx"
)

// QueryItem retrieves a row from the database
func (g *gateway) QueryItem(ctx context.Context, item interface{}, qry string, args ...interface{}) error {

	// Validate 'item' is a pointer to a struct
	vItem, err := isStructPtr(item)
	if err != nil {
		return err
	}

	// Open connection with tenant isolation
	conn, err := g.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Query the rows
	rows, err := conn.QueryContext(ctx, qry, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Create a slice of the struct
	vListPtr := reflect.New(reflect.SliceOf(vItem.Type()))

	// Scan rows
	if err := sqlx.StructScan(rows, vListPtr.Interface()); err != nil {
		return err
	}

	if vListPtr.Elem().Len() < 1 {
		return sql.ErrNoRows
	}

	vItem.Set(vListPtr.Elem().Index(0))
	return nil
}
