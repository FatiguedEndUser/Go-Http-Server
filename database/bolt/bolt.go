// Package bolt provides a BoltDB implementation of the database.Database interface.
// Users are stored as JSON blobs keyed by Name in a single BoltDB bucket.
package bolt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/boltdb/bolt"
	"Http-Server/database"
)

const (
	dbName     = "test.db"
	bucketName = "users"
)

// Bolt satisfies the database.Database interface using BoltDB as the persistence backend.
type Bolt struct {
	db *bolt.DB
}

// New creates a new BoltDB instance at the specified directory.
// Ensures the users bucket exists, creating it if necessary.
// Returns an error if the database cannot be opened or bucket creation fails.
func New(ctx context.Context, directory string) (*Bolt, error) {
	path := filepath.Join(directory, dbName)

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open database at %s: %w", path, err)
	}

	// Create bucket if it doesn't exist
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	return &Bolt{
		db: db,
	}, nil
}

// Close releases the database connection. Always call when finished.
func (b *Bolt) Close(ctx context.Context) error {
	return b.db.Close()
}

// Create inserts a new user into the database.
// Returns an error if marshaling or persistence fails.
func (b *Bolt) Create(ctx context.Context, user database.User) error {
	if user.Name == "" {
		return fmt.Errorf("user name cannot be empty")
	}

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	// Execute write transaction and propagate any errors
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		return bucket.Put([]byte(user.Name), data)
	})
	if err != nil {
		return fmt.Errorf("failed to persist user %q: %w", user.Name, err)
	}

	return nil
}

// Get retrieves a user by name.
// Returns nil if the user does not exist (not an error condition).
// On JSON decode failure, logs a warning and returns nil instead of crashing.
func (b *Bolt) Get(ctx context.Context, name string) *database.User {
	if name == "" {
		return nil
	}

	var user *database.User
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return nil
		}

		data := bucket.Get([]byte(name))
		if len(data) == 0 {
			return nil // User not found
		}

		var u database.User
		if err := json.Unmarshal(data, &u); err != nil {
			// Changed from log.Fatalf to log.Printf - don't crash on corrupt data
			log.Printf("warning: failed to decode user %q: %v", name, err)
			return nil
		}
		user = &u
		return nil
	})

	if err != nil {
		log.Printf("error reading user %q: %v", name, err)
		return nil
	}

	return user
}

// Update modifies an existing user's information.
// Fetches current data, applies changes, and persists the result.
// Returns the updated user on success, or an error if the user doesn't exist or update fails.
func (b *Bolt) Update(ctx context.Context, input database.User) (*database.User, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("user name cannot be empty")
	}

	// Read existing data first
	var raw []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		raw = bucket.Get([]byte(input.Name))
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(raw) == 0 {
		return nil, fmt.Errorf("user %q not found", input.Name)
	}

	var current database.User
	if err := json.Unmarshal(raw, &current); err != nil {
		return nil, fmt.Errorf("failed to decode existing user: %w", err)
	}

	// Apply updates (allows partial updates)
	current.Email = input.Email
	current.Age = input.Age
	// Name stays the same as the lookup key

	data, err := json.Marshal(current)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated user: %w", err)
	}

	// Persist changes
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		return bucket.Put([]byte(input.Name), data)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to persist changes: %w", err)
	}

	return &current, nil
}

// Delete removes a user from the database by name.
// Returns an error if deletion fails.
func (b *Bolt) Delete(ctx context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		return bucket.Delete([]byte(name))
	})
	if err != nil {
		return fmt.Errorf("failed to delete user %q: %w", name, err)
	}

	return nil
}