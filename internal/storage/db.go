package storage

import (
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/Oguzyildirim/go-counter/internal"
)

//Driver drive the db
type Driver struct {
	dir   string
	mutex *sync.Mutex
}

// Creates new driver
func New(dir string) *Driver {
	driver := &Driver{
		dir:   dir,
		mutex: &sync.Mutex{},
	}
	return driver
}

// Insert creates a new record at db
func (d *Driver) Insert() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	now := time.Now()
	formatted := now.Format(time.RFC1123)
	data := formatted + "\n"

	f, err := os.OpenFile(d.dir, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "OpenFile failed")
	}
	defer f.Close()

	if _, err := f.WriteString(data); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "WriteString failed")
	}
	return nil
}

// Get finds a new record
func (d *Driver) Get() (string, error) {
	b, err := ioutil.ReadFile(d.dir)
	if err != nil {
		return "", internal.WrapErrorf(err, internal.ErrorCodeUnknown, "ioutil read file failed")
	}
	return string(b), nil
}
