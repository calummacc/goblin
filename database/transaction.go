// goblin/database/transaction.go
// Package database provides database access functionality for the Goblin Framework.
// The transaction module implements database transaction management to ensure
// data consistency and provide atomic operations.
package database

import (
	"context"

	"gorm.io/gorm"
)

// TransactionManager manages database transactions in the Goblin Framework.
// It provides a clean interface for executing operations within transactions,
// handling commit and rollback automatically based on the success or failure
// of the operations.
type TransactionManager struct {
	// orm is the database ORM instance used for creating transactions
	orm *ORM
}

// NewTransactionManager creates a new transaction manager instance.
// It initializes the manager with the provided ORM instance.
//
// Parameters:
//   - orm: The ORM instance to use for database operations
//
// Returns:
//   - *TransactionManager: A new transaction manager instance
func NewTransactionManager(orm *ORM) *TransactionManager {
	return &TransactionManager{
		orm: orm,
	}
}

// Transaction runs a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
//
// Example:
//
//	tm.Transaction(ctx, func(tx *gorm.DB) error {
//	    // Create a user
//	    if err := tx.Create(&user).Error; err != nil {
//	        return err
//	    }
//	    // Create a profile for the user
//	    if err := tx.Create(&profile).Error; err != nil {
//	        return err
//	    }
//	    return nil
//	})
//
// Parameters:
//   - ctx: The context for the transaction, which can be used for timeouts or cancellation
//   - fn: The function to execute within the transaction, receiving the transaction handle
//
// Returns:
//   - error: Any error that occurred during the transaction
func (tm *TransactionManager) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return tm.orm.DB.WithContext(ctx).Transaction(fn)
}

// Transactional is a decorator for running a function within a transaction.
// It can be used to wrap services or repositories to make their operations transactional.
// This provides a more object-oriented approach to transaction management compared
// to directly using the TransactionManager.
type Transactional struct {
	// manager is the transaction manager to use for transaction handling
	manager *TransactionManager
}

// NewTransactional creates a new transactional decorator instance.
// It initializes the decorator with the provided transaction manager.
//
// Parameters:
//   - manager: The transaction manager to use
//
// Returns:
//   - *Transactional: A new transactional decorator instance
func NewTransactional(manager *TransactionManager) *Transactional {
	return &Transactional{
		manager: manager,
	}
}

// Run executes a function within a transaction.
// This is a convenience method that delegates to the transaction manager.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
//
// Parameters:
//   - ctx: The context for the transaction
//   - fn: The function to execute within the transaction
//
// Returns:
//   - error: Any error that occurred during the transaction
func (t *Transactional) Run(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return t.manager.Transaction(ctx, fn)
}
