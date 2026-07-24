package database

import (
	"context"
)

//This represents the operations done on the database
type Database interface {
	Create(ctx context.Context, data []byte) error
	
}