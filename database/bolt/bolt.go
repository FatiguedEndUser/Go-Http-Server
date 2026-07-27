package bolt

import (
	"Http-Server/database"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const (
	dbName     = "test.db"
	bucketName = "users"
)

// Bolt is the bolt database
// It satisfies the Database interface
type Bolt struct {
	db *bolt.DB
}

// New returns a new Bolt implementation
func New(ctx context.Context, directory string) (*Bolt, error) {
	db, err := bolt.Open(fmt.Sprintf("%s/%s", directory, dbName), 0600, nil)
	if err != nil {
		return nil, err
	}

	//Ensure that the bucket exists, if not, create it
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &Bolt{
		db: db,
	}, nil
}

type UserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// Closes the bolt database
func (b *Bolt) Close(ctx context.Context) error {
	return b.db.Close()
}

// Create implements the database interface
func (b *Bolt) Create(ctx context.Context, user database.User) error {
	userinfo := UserInfo{
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}

	v, err := json.Marshal(userinfo)
	if err != nil {
		return err
	}
	b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.Put([]byte(user.Name), v)
	})
	return nil
}

// Get implements the database interface
func (b *Bolt) Get(ctx context.Context, name string) (user *database.User) {
	var raw []byte
	b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		raw = b.Get([]byte(name))
		return nil
	})
	if len(raw) == 0 {
		return nil
	}

	var u database.User
	err := json.Unmarshal(raw, &u)
	if err != nil {
		log.Fatalf("Database Corruption %w", err)
	}
	user = &u
	return
}

func (b *Bolt) Update(ctx context.Context, input database.User) (*database.User, error) {
	var raw []byte
	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		raw = bucket.Get([]byte(input.Name))
		return nil
	})

	var current database.User
	err := json.Unmarshal(raw, &current)
	if err != nil {
		return nil, fmt.Errorf("Database Corruption %w", err)
	}
	current.Name = input.Name
	current.Age = input.Age
	current.Email = input.Email

	v, err := json.Marshal(current)
	if err != nil {
		return nil, fmt.Errorf("Could not marshal user: %w", err)
	}
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		return bucket.Put([]byte(input.Name), v)
	})

	if err != nil {
		return nil, err
	}

	return &current, nil
}

func (b *Bolt) Delete(ctx context.Context, name string) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		return bucket.Delete([]byte(name))
	})
	return err
}