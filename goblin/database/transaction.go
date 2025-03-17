// goblin/database/transaction.go
package database

import (
	"context"

	"gorm.io/gorm"
)

// TransactionManager manages database transactions
type TransactionManager struct {
	orm *ORM
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(orm *ORM) *TransactionManager {
	return &TransactionManager{
		orm: orm,
	}
}

// Transaction runs a function within a transaction
func (tm *TransactionManager) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return tm.orm.DB.WithContext(ctx).Transaction(fn)
}

// Transactional is a decorator for running a function within a transaction
type Transactional struct {
	manager *TransactionManager
}

// NewTransactional creates a new transactional decorator
func NewTransactional(manager *TransactionManager) *Transactional {
	return &Transactional{
		manager: manager,
	}
}

// Run runs a function within a transaction
func (t *Transactional) Run(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return t.manager.Transaction(ctx, fn)
}
