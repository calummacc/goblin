// goblin/database/orm.go
// Package database provides database access functionality for the Goblin Framework.
// It implements an ORM layer based on GORM, providing a clean interface for
// database operations and model management.
package database

import (
	"gorm.io/gorm"
)

// Config defines the configuration options for the ORM.
// It contains settings needed to establish a database connection.
type Config struct {
	// DSN (Data Source Name) specifies the connection string for the database.
	// Format depends on the database driver being used.
	// Examples:
	// - MySQL: "user:password@tcp(host:port)/dbname?charset=utf8&parseTime=True&loc=Local"
	// - PostgreSQL: "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable"
	// - SQLite: "file:test.db?cache=shared"
	DSN string
}

// ORM represents a database ORM instance.
// It wraps the GORM DB instance and provides additional functionality
// specific to the Goblin Framework, making database operations more
// convenient and consistent.
type ORM struct {
	// DB is the underlying GORM database instance
	DB *gorm.DB
}

// NewORM creates a new ORM instance.
// This function initializes the database connection and returns an ORM wrapper.
// In the current implementation, it returns a mock ORM for development purposes,
// but in a production environment it would establish a real database connection.
//
// In the future, it should accept a Config parameter to customize the connection.
//
// Returns:
//   - *ORM: A new ORM instance
//   - error: Any error that occurred during initialization
func NewORM() (*ORM, error) {
	// TODO: Implement real database connection
	// Temporarily return a mock ORM
	return &ORM{}, nil
}

// AutoMigrate automatically migrates the database schema for the provided models.
// It creates or updates database tables to match the model definitions.
// This should typically be called during application startup.
//
// Example:
//
//	orm.AutoMigrate(&User{}, &Product{}, &Order{})
//
// Parameters:
//   - models: A variadic list of model types to migrate
//
// Returns:
//   - error: Any error that occurred during migration
func (o *ORM) AutoMigrate(models ...interface{}) error {
	// TODO: Implement real migration when DB is available
	// Temporarily do nothing
	return nil
}

// Repository provides a base implementation for database access.
// It encapsulates common database operations and provides a clean interface
// for working with database entities. It's designed to be embedded in
// domain-specific repositories.
//
// Example:
//
//	type UserRepository struct {
//	    *database.Repository
//	}
//
//	func (r *UserRepository) FindByEmail(email string) (*User, error) {
//	    var user User
//	    if err := r.ORM.DB.Where("email = ?", email).First(&user).Error; err != nil {
//	        return nil, err
//	    }
//	    return &user, nil
//	}
type Repository struct {
	// ORM is the database ORM instance used by this repository
	ORM *ORM
}

// NewRepository creates a new repository instance.
// It initializes the repository with the provided ORM instance.
// This function is typically used with dependency injection to provide
// repositories to services and controllers.
//
// Parameters:
//   - orm: The ORM instance to use for database operations
//
// Returns:
//   - *Repository: A new repository instance
func NewRepository(orm *ORM) *Repository {
	return &Repository{
		ORM: orm,
	}
}
