package database

import (
	"context"
)

//User information Types
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

//This represents the operations done on the database
type Database interface {
	Create(ctx context.Context, user User) error
	Get(ctx context.Context, name string) *User
}