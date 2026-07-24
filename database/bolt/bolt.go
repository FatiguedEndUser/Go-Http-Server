package bolt

import (
	"context"
	"fmt"
)

//Bolt is the bolt database
//It satisfies the Database interface
type Bolt struct {
	
}

//New returns a new Bolt implementation
func New(ctx context.Context) *Bolt {
	return &Bolt{}
}

//Create implements the database interface
func (b *Bolt) Create(ctx context.Context, data []byte) error {
	fmt.Println(data)
	return nil
}