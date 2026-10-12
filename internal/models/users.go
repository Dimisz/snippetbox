package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// UserModel struct wraps a database connection pool
type UserModel struct {
	DB *sql.DB
}

// Insert() method is used to add a new record to the 'users' table
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate() method is used to verify if a user with the provided email address and password exists, returns user's id if the user exists, or error if they do not
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Exists() is used to check if a user with a specific ID exists
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
