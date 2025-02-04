package data

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Manager represents the manager to manage the data consistency
type Manager struct {
	db *sqlx.DB
}

// RunInTransaction runs the f with the transaction queryable inside the context
func (m *Manager) RunInTransaction(ctx context.Context, f func(tctx context.Context) error) error {
	tx, err := m.db.Beginx()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error when creating transction: %v", err)
	}

	ctx = NewContext(ctx, tx)
	err = f(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error when committing transaction: %v", err)
	}

	return nil
}

// NewManager creates a new manager
func NewManager(db *sqlx.DB) *Manager {
	return &Manager{
		db: db,
	}
}
