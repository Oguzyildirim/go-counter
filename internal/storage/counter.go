package storage

import (
	"github.com/Oguzyildirim/go-counter/internal"
)

// Counter represents the repository used for interacting with Counter records
type Counter struct {
	db *Driver
}

// NewCounter instantiates the Counter repository
func NewCounter(dir string) *Counter {
	return &Counter{
		db: New(dir),
	}
}

// Create inserts a new ıser record
func (c *Counter) Create() error {
	err := c.db.Insert()
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "create failed")
	}
	return nil
}

// Find
func (c *Counter) Find() (string, error) {
	result, err := c.db.Get()
	if err != nil {
		return "", internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "find failed")
	}

	return result, nil
}
