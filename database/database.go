// Package database defines the core data structures and abstractions for user storage.
// It provides a contract (Database interface) that all persistence implementations must satisfy,
// enabling loose coupling between the HTTP layer and storage backends.
package database

import (
	"context"
)

// User represents a person stored in the database.
// All fields are required and validated at the application layer.
type User struct {
	// Name is the unique identifier for the user (must be non-empty).
	Name string `json:"name"`
	
	// Email is the user's email address (not enforced as unique in this implementation).
	Email string `json:"email"`
	
	// Age is the user's age in years (application-level validation required).
	Age int `json:"age"`
}

// Database defines the contract for persistent user storage operations.
// Implementations must be concurrency-safe and respect context cancellation.
//
// Methods follow standard Go conventions:
//   - Create/Update/Delete return errors on failure
//   - Get returns nil when no record exists (not an error)
//   - All methods accept context for timeout/cancellation control
type Database interface {
	// Create adds a new user to the database.
	// Returns an error if the user already exists or the operation fails.
	Create(ctx context.Context, user User) error
	
	// Get retrieves a user by name.
	// Returns nil if no user with the given name exists.
	Get(ctx context.Context, name string) *User
	
	// Update modifies an existing user's information.
	// Returns the updated user on success, or an error if the user doesn't exist.
	Update(ctx context.Context, user User) (*User, error)
	
	// Delete removes a user from the database by name.
	// Returns an error if the user doesn't exist or the operation fails.
	Delete(ctx context.Context, name string) error
}