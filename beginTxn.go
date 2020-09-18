package pgdb

import (
	"context"
	"fmt"
)

// BeginTxn begins a new database transaction
func (g *gateway) BeginTxn(ctx context.Context) (Txn, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	if tenantID, ok := ctx.Value(TenantID).(string); ok {
		cmd := fmt.Sprintf("SET app.tenant_id = '%v';", tenantID)
		if _, err := tx.ExecContext(ctx, cmd); err != nil {
			return nil, err
		}
	}
	return &txn{
		tx: tx,
	}, nil
}
