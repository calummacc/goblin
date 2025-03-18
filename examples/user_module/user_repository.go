// goblin/examples/user_module/user_repository.go
package user_module

import (
	"errors"
	"goblin/database"
)

// UserRepository handles user data access
type UserRepository struct {
	*database.Repository
}

// NewUserRepository creates a new user repository
func NewUserRepository(repo *database.Repository) *UserRepository {
	return &UserRepository{
		Repository: repo,
	}
}

// FindAll returns all users
func (r *UserRepository) FindAll() ([]User, error) {
	var users []User
	result := r.ORM.DB.Find(&users)
	return users, result.Error
}

// FindByID returns a user by ID
func (r *UserRepository) FindByID(id uint) (User, error) {
	var user User
	result := r.ORM.DB.First(&user, id)
	if result.RowsAffected == 0 {
		return User{}, errors.New("user not found")
	}
	return user, result.Error
}

// Create creates a new user
func (r *UserRepository) Create(user User) (User, error) {
	result := r.ORM.DB.Create(&user)
	return user, result.Error
}

// Update updates a user
func (r *UserRepository) Update(id uint, user User) (User, error) {
	var existingUser User
	result := r.ORM.DB.First(&existingUser, id)
	if result.RowsAffected == 0 {
		return User{}, errors.New("user not found")
	}

	user.ID = id
	result = r.ORM.DB.Save(&user)
	return user, result.Error
}

// Delete deletes a user
func (r *UserRepository) Delete(id uint) error {
	result := r.ORM.DB.Delete(&User{}, id)
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return result.Error
}
