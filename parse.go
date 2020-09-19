package sqldb

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/lib/pq"
)

func (g *gateway) parseRow(row *sql.Row, item interface{}) error {

	// Validate 'item' is a pointer to a struct
	vItem, err := isStructPtr(item)
	if err != nil {
		return err
	}

	// Scan the SQL row into the item
	return g.scan(row, vItem)
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
		if err := g.scan(rows, vItemPtr.Elem()); err != nil {
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

func (g *gateway) scan(s scannable, vItem reflect.Value) error {

	// Create slice of pointers
	// Each pointer's value type matches a field type of the struct being written to
	ptrs := make([]interface{}, vItem.NumField())
	driverPtrs := make([]interface{}, vItem.NumField())
	for i := 0; i < vItem.NumField(); i++ {
		vField := vItem.Field(i)
		ptr := reflect.New(vField.Type()).Interface()
		ptrs[i] = ptr
		driverPtrs[i] = driverCast(ptr, vField, g.driver)
	}

	// Scan row in the pointer values
	if err := s.Scan(driverPtrs...); err != nil {
		return err
	}

	// Assign the pointer values to struct fields
	for i := 0; i < len(ptrs); i++ {
		vItem.Field(i).Set(reflect.ValueOf(ptrs[i]).Elem())
	}

	return nil
}

func driverCast(ptr interface{}, v reflect.Value, driver string) interface{} {
	switch driver {
	case "postgres":
		return postgresCast(ptr, v)
	case "sqlite3":
		return sqliteCast(ptr, v)
	default:
		panic(fmt.Errorf("Driver '%v' unsupported", driver))
	}
}

func postgresCast(ptr interface{}, v reflect.Value) interface{} {
	switch v.Kind() {
	case reflect.Slice:
		switch v.Type().Name() {
		case "[]byte":
			return ptr
		default:
			return pq.Array(ptr)
		}

	default:
		return ptr
	}
}

func sqliteCast(ptr interface{}, v reflect.Value) interface{} {
	return ptr
}
