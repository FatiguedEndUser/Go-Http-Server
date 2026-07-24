package bolt

import (
	"context"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const (
	dbName = "test.db"
	bucketName = "users"
)


//Bolt is the bolt database
//It satisfies the Database interface
type Bolt struct {
	db *bolt.DB
}

//New returns a new Bolt implementation
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

//Closes the bolt database
func (b *Bolt) Close(ctx context.Context) error {
	return b.db.Close()
}

//Create implements the database interface
func (b *Bolt) Create(ctx context.Context, data []byte) error {
	fmt.Println(string(data))
	return nil
}