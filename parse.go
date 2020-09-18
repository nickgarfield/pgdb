package sqldb

import (
	"database/sql"
	"reflect"

	"github.com/lib/pq"
)

func parseRow(row *sql.Row, item interface{}) error {

	// Validate 'item' is a pointer to a struct
	vItem, err := isStructPtr(item)
	if err != nil {
		return err
	}

	// Scan the SQL row into the item
	return scan(row, vItem)
}

func parseRows(rows *sql.Rows, list interface{}) error {

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
		if err := scan(rows, vItemPtr.Elem()); err != nil {
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

func scan(s scannable, vItem reflect.Value) error {

	// Create slice of pointers
	// Each pointer's value type matches a field type of the struct being written to
	ptrs := make([]interface{}, vItem.NumField())
	pqWrappedPtrs := make([]interface{}, vItem.NumField())
	for i := 0; i < vItem.NumField(); i++ {
		vField := vItem.Field(i)
		ptr := reflect.New(vField.Type()).Interface()
		ptrs[i] = ptr
		switch vField.Kind() {
		case reflect.Slice:
			pqWrappedPtrs[i] = pq.Array(ptr)
		default:
			pqWrappedPtrs[i] = ptr
		}
	}

	// Scan row in the pointer values
	if err := s.Scan(pqWrappedPtrs...); err != nil {
		return err
	}

	// Assign the pointer values to struct fields
	for i := 0; i < len(ptrs); i++ {
		vItem.Field(i).Set(reflect.ValueOf(ptrs[i]).Elem())
	}

	return nil
}
