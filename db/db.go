package db

import (
	"errors"
	"sync"
)

// Database is the interface that defines the in-memory database operations
type Database interface {
	Get(string) (string, error)
	Set(string, string)
}

// Database is the in-memory database object that holds the data
type database struct {
	M    *sync.RWMutex
	Data map[string]string
}

// Get returns the value stored in the map[key] or an error
func (d *database) Get(key string) (string, error) {
	d.M.RLock()
	val, ok := d.Data[key]
	d.M.RUnlock()
	if !ok {
		return "", errors.New("Key does not exist")
	}
	return val, nil
}

// Set sets the value of map[key] to value
func (d *database) Set(key string, val string) {
	d.M.Lock()
	defer d.M.Unlock()
	d.Data[key] = val
}

// NewDatabase initializes and returns the in memory database objcet
func NewDatabase() Database {
	return &database{
		M:    &sync.RWMutex{},
		Data: make(map[string]string),
	}
}
