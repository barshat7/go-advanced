package localstore

import (
	"errors"
	"sync"
)

var store = struct {
	sync.RWMutex
	m map[string] string
}{m: make(map[string] string)}

// ErrorNoSuchKey just err
var ErrorNoSuchKey = errors.New("No Such Key")

// Put put
func Put(key string, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()
	return nil
}

// Get get
func Get(key string) (string, error) {
	store.RLock()	
	value, ok := store.m[key]
	store.RUnlock()
	if !ok {
		return "", ErrorNoSuchKey
	}
	return value , nil
}

// Delete delete
func Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()
	return nil
}
