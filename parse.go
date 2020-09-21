package sqldb

import (
	"database/sql"
	"reflect"

	"github.com/jmoiron/sqlx"
)

func (g *gateway) parseRow(rows *sql.Rows, item interface{}) error {

	// Validate 'item' is a pointer to a struct
	_, err := isStructPtr(item)
	if err != nil {
		return err
	}

	// Scan the SQL row into the item
	rows.Next()
	return sqlx.StructScan(rows, item)
}

func (g *gateway) parseRows(rows *sql.Rows, list interface{}) error {

	// Validate 'list' is a pointer to a slice
	vList, err := isSlicePtr(list)
	if err != nil {
		return err
	}

	// Get the item type of the slice
	tItem := vList.Type().Elem()

	// For each SQL row...
	for rows.Next() {

		// Generate a pointer to a new item
		vItemPtr := reflect.New(tItem)

		// Scan the SQL row into the item
		if err := sqlx.StructScan(rows, vItemPtr.Interface()); err != nil {
			return err
		}

		// Append the item to the list
		vList.Set(reflect.Append(vList, vItemPtr.Elem()))
	}

	return nil
}

type scannable interface {
	Scan(dest ...interface{}) error
}
