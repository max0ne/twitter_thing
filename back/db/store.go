package db

import (
	"fmt"
	"strconv"
	"sync"
)

// Store - -
type Store struct {
	m     map[string]string
	mLock *sync.Mutex
}

// GetArgs - -
type GetArgs struct {
	Key string
}

// PutArgs - -
type PutArgs struct {
	Key string
	Val string
}

// DelArgs - -
type DelArgs struct {
	Key string
}

// GetReply - -
type GetReply struct {
	Val string
}

// HasReply - -
type HasReply struct {
	Has bool
}

// IncIDArgs - -
type IncIDArgs struct {
	TableName string
}

// IncIDReply - -
type IncIDReply struct {
	ID string
}

// GetMReply - -
type GetMReply struct {
	M map[string]string
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
func (s *Store) GetM(args interface{}, reply *GetMReply) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()

	*reply = GetMReply{
		M: s.m,
	}
	return nil
}

// Get - -
func (s *Store) Get(args GetArgs, reply *GetReply) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()

	*reply = GetReply{Val: s.m[args.Key]}
	return nil
}

// Put - -
func (s *Store) Put(args PutArgs, ack *bool) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()

	s.m[args.Key] = args.Val
	return nil
}

// Del - -
func (s *Store) Del(args DelArgs, ack *bool) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()

	delete(s.m, args.Key)
	return nil
}

// Has - -
func (s *Store) Has(args GetArgs, reply *HasReply) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()

	_, ok := s.m[args.Key]
	*reply = HasReply{Has: ok}
	return nil
}

// IncID - -
func (s *Store) IncID(args IncIDArgs, reply *IncIDReply) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()

	key := "$IncID_" + args.TableName
	if s.m[key] == "" {
		s.m[key] = "1"
	} else {
		id, err := strconv.ParseInt(s.m[key], 10, 64)
		if err != nil {
			return err
		}
		s.m[key] = fmt.Sprintf("%d", id+1)
	}
	*reply = IncIDReply{ID: s.m[key]}
	return nil
}
