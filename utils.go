package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

// TenantID is the TENANT_ID key in the Go context
const TenantID = "TENANT_ID"

// connect opens a new database connection with tenant isolation
func (g *gateway) connect(ctx context.Context) (*sql.Conn, error) {
	conn, err := g.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	if tenantID, ok := ctx.Value(TenantID).(string); ok {
		cmd := fmt.Sprintf("SET app.tenant_id = '%v';", tenantID)
		if _, err = conn.ExecContext(ctx, cmd); err != nil {
			conn.Close()
			return nil, err
		}
	}
	return conn, nil
}

// isStructPtr validates 'item' is a pointer to a struct
func isStructPtr(item interface{}) (reflect.Value, error) {
	vItemPtr := reflect.ValueOf(item)
	if vItemPtr.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("'item' must be a pointer to a struct")
	}
	vItem := vItemPtr.Elem()
	if vItem.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("'item' must be a pointer to a struct")
	}
	return vItem, nil
}

// isSlicePtr validates 'list' is a pointer to a slice
func isSlicePtr(list interface{}) (reflect.Value, error) {
	vListPtr := reflect.ValueOf(list)
	if vListPtr.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("'list' must be a pointer to a slice")
	}
	vList := vListPtr.Elem()
	if vList.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("'list' must be a pointer to a slice")
	}
	return vList, nil
}
