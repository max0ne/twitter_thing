package db

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
)

// Store - -
type Store struct {
	m     map[string]string
	mLock *sync.Mutex
}

// Table - -
type Table struct {
	Name  string
	store *Store

	incrementID int64
}

// NewStore - -
func NewStore() *Store {
	store := Store{
		m:     map[string]string{},
		mLock: &sync.Mutex{},
	}
	return &store
}

// GetM - for debug only
func (s *Store) GetM() map[string]string {
	return s.m
}

// NewTable - -
func (s *Store) NewTable(name string) *Table {
	return &Table{
		Name:  name,
		store: s,
	}
}

// Put - -
func (t *Table) Put(key, val string) {
	t.store.mLock.Lock()
	defer t.store.mLock.Unlock()

	t.store.m[t.key(key)] = val
}

// PutObj - -
func (t *Table) PutObj(key string, val interface{}) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	t.Put(key, string(bytes))
	return nil
}

// Del - -
func (t *Table) Del(key string) {
	t.store.mLock.Lock()
	defer t.store.mLock.Unlock()

	delete(t.store.m, t.key(key))
}

// Get - -
func (t *Table) Get(key string) string {
	t.store.mLock.Lock()
	defer t.store.mLock.Unlock()

	return t.store.m[t.key(key)]
}

// GetObj - -
func (t *Table) GetObj(key string, target interface{}) error {
	jsonString := t.Get(key)
	if jsonString == "" {
		return nil
	}
	return json.Unmarshal([]byte(jsonString), target)
}

// Has has key
func (t *Table) Has(key string) bool {
	t.store.mLock.Lock()
	defer t.store.mLock.Unlock()

	_, ok := t.store.m[t.key(key)]
	return ok
}

func (t *Table) key(key string) string {
	return t.Name + "_" + key
}

// IncID get auto increment id
func (t *Table) IncID() string {
	return fmt.Sprintf("%d", atomic.AddInt64(&t.incrementID, 1))
}
