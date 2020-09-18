package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/lib/pq"
)

// Writable provides an interface for values that are writable to the database
type Writable interface {
	DatabaseTable() string
}

// Txn provides an interface for database transactions
type Txn interface {
	Insert(ctx context.Context, item Writable, clause string) error
	Rollback() error
	Commit() error
}

type txn struct {
	tx *sql.Tx
}

// Insert appends an SQL INSERT command for a writeable value to a database transaction
func (txn *txn) Insert(ctx context.Context, item Writable, clause string) error {

	// Validate 'item' is a pointer to a struct
	vItem, err := isStructPtr(item)
	if err != nil {
		return err
	}
	tItem := vItem.Type()

	// Prepare arguments for SQL command
	fieldNames := []string{}
	fieldPlaceholders := []string{}
	fieldValues := []interface{}{}

	// For each field in the struct...
	for i := 0; i < vItem.NumField(); i++ {

		// Extract the field name in the json tag
		if tag := tItem.Field(i).Tag.Get("json"); tag != "" {

			// Save the field name and placeholder
			fieldNames = append(fieldNames, tag)
			fieldPlaceholders = append(fieldPlaceholders, fmt.Sprintf("$%v", len(fieldPlaceholders)+1))

			// Save the field value
			switch vItem.Field(i).Kind() {
			case reflect.Slice:
				fieldValues = append(fieldValues, pq.Array(vItem.Field(i).Interface()))
			default:
				fieldValues = append(fieldValues, vItem.Field(i).Interface())
			}
		}
	}

	// Build and execute the command
	cmd := fmt.Sprintf(`INSERT INTO %v (%v) VALUES (%v) %v`,
		item.DatabaseTable(),
		strings.Join(fieldNames, ", "),
		strings.Join(fieldPlaceholders, ", "),
		clause,
	)
	_, err = txn.tx.ExecContext(ctx, cmd, fieldValues...)
	if err != nil {
		txn.Rollback()
	}
	return err
}

// Rollback rolls back a database transaction
func (txn *txn) Rollback() error {
	return txn.tx.Rollback()
}

// Commit commits a database transaction
func (txn *txn) Commit() error {
	return txn.tx.Commit()
}
