// goblin/database/orm.go
package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ORM represents a database ORM
type ORM struct {
	DB *gorm.DB
}

// Config represents database configuration
type Config struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Options  map[string]string
}

// NewORM creates a new ORM instance
func NewORM(config Config) (*ORM, error) {
	var db *gorm.DB
	var err error

	switch config.Driver {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.Database), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &ORM{
		DB: db,
	}, nil
}

// AutoMigrate runs auto migration for the provided models
func (o *ORM) AutoMigrate(models ...interface{}) error {
	return o.DB.AutoMigrate(models...)
}

// Repository is a base repository for database access
type Repository struct {
	ORM *ORM
}

// NewRepository creates a new repository
func NewRepository(orm *ORM) *Repository {
	return &Repository{
		ORM: orm,
	}
}
