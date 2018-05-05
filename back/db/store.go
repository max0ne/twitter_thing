package db

import (
	"encoding/gob"
	"fmt"
	"sync"
)

// Store - -
type Store struct {
	m      map[string]string
	incIDs map[string]int
	mLock  *sync.Mutex

	emitCommand func(cmd interface{}) error
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
func NewStore(emitCommand func(cmd interface{}) error) *Store {
	gob.Register(PutArgs{})
	gob.Register(DelArgs{})
	gob.Register(IncIDArgs{})
	store := Store{
		m:           map[string]string{},
		incIDs:      map[string]int{},
		mLock:       &sync.Mutex{},
		emitCommand: emitCommand,
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

// Has - -
func (s *Store) Has(args GetArgs, reply *HasReply) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()

	_, ok := s.m[args.Key]
	*reply = HasReply{Has: ok}
	return nil
}

// Put - -
func (s *Store) Put(args PutArgs, ack *bool) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.emitCommand(args)
}

// Del - -
func (s *Store) Del(args DelArgs, ack *bool) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.emitCommand(args)
}

// IncID - -
func (s *Store) IncID(args IncIDArgs, reply *IncIDReply) error {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.emitCommand(args)
}

// ---
// internal write methods
// ---
func (s *Store) put(args PutArgs) {
	fmt.Println("put", args, s.m, s.m == nil)
	s.m[args.Key] = args.Val
}
func (s *Store) del(args DelArgs) {
	fmt.Println("del", args)
	delete(s.m, args.Key)
}
func (s *Store) incid(args IncIDArgs) IncIDReply {
	fmt.Println("incid", args)
	s.incIDs[args.TableName]++
	return IncIDReply{
		ID: fmt.Sprintf("%d", s.incIDs[args.TableName]),
	}
}

func (s *Store) processWriteCommand(cmd interface{}) {
	s.mLock.Lock()
	defer s.mLock.Unlock()

	fmt.Println("process cmd", cmd)

	putArg, ok := cmd.(PutArgs)
	if ok {
		s.put(putArg)
		return
	}

	delArg, ok := cmd.(DelArgs)
	if ok {
		s.del(delArg)
		return
	}

	incIDArg, ok := cmd.(IncIDArgs)
	if ok {
		s.incid(incIDArg)
		return
	}
}
