package db

import (
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/max0ne/twitter_thing/back/vr"
)

// Store - -
type Store struct {
	m      map[string]string
	incIDs map[string]int
	mLock  *sync.Mutex

	emitCommand func(cmd vr.Command) error
	server      *Server
}

type Args struct {
	Kind string
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

// IncIDArgs - -
type IncIDArgs struct {
	TableName string
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

// IncIDReply - -
type IncIDReply struct {
	ID string
}

// GetMReply - -
type GetMReply struct {
	M map[string]string
}

// NewStore - -
func NewStore(emitCommand func(cmd vr.Command) error, server *Server) *Store {
	gob.Register(PutArgs{})
	gob.Register(DelArgs{})
	gob.Register(IncIDArgs{})
	store := Store{
		mLock:       &sync.Mutex{},
		emitCommand: emitCommand,
		server:      server,
	}
	store.reset()
	return &store
}

// GetM - for debug only
func (s *Store) GetM(args interface{}, reply *GetMReply) error {
	// get from primary if not me
	if !vr.IsPrimary(s.server.vrServer) {
		val, err := s.server.dbPeerClients[vr.Primary(s.server.vrServer)].GetM()
		if err != nil {
			return err
		}
		*reply = GetMReply{
			M: val,
		}
		return nil
	}

	s.mLock.Lock()
	defer s.mLock.Unlock()

	*reply = GetMReply{
		M: s.m,
	}
	return nil
}

// Get - -
func (s *Store) Get(args GetArgs, reply *GetReply) error {
	// get from primary if not me
	if !vr.IsPrimary(s.server.vrServer) {
		fmt.Println("getting from primary")
		val, err := s.server.dbPeerClients[vr.Primary(s.server.vrServer)].Get(args.Key)
		if err != nil {
			return err
		}
		*reply = GetReply{
			Val: val,
		}
		return nil
	}

	s.mLock.Lock()
	defer s.mLock.Unlock()

	*reply = GetReply{Val: s.m[args.Key]}
	return nil
}

// Has - -
func (s *Store) Has(args GetArgs, reply *HasReply) error {
	// get from primary if not me
	if !vr.IsPrimary(s.server.vrServer) {
		val, err := s.server.dbPeerClients[vr.Primary(s.server.vrServer)].Has(args.Key)
		if err != nil {
			return err
		}
		*reply = HasReply{
			Has: val,
		}
		return nil
	}

	s.mLock.Lock()
	defer s.mLock.Unlock()

	_, ok := s.m[args.Key]
	*reply = HasReply{Has: ok}
	return nil
}

// Put - -
func (s *Store) Put(args PutArgs, ack *bool) error {
	// get from primary if not me
	if !vr.IsPrimary(s.server.vrServer) {
		fmt.Println(s.server.dbPeerClients, vr.Primary(s.server.vrServer))
		err := s.server.dbPeerClients[vr.Primary(s.server.vrServer)].Put(args.Key, args.Val)
		if err != nil {
			return err
		}
		return nil
	}

	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.emitCommand(vr.Command{
		Kind:  "Put",
		Value: args,
	})
}

// Del - -
func (s *Store) Del(args DelArgs, ack *bool) error {
	// get from primary if not me
	if !vr.IsPrimary(s.server.vrServer) {
		err := s.server.dbPeerClients[vr.Primary(s.server.vrServer)].Del(args.Key)
		if err != nil {
			return err
		}
		return nil
	}

	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.emitCommand(vr.Command{
		Kind:  "Del",
		Value: args,
	})
}

// IncID - -
func (s *Store) IncID(args IncIDArgs, reply *IncIDReply) error {
	// get from primary if not me
	if !vr.IsPrimary(s.server.vrServer) {
		val, err := s.server.dbPeerClients[vr.Primary(s.server.vrServer)].IncID(args.TableName)
		if err != nil {
			return err
		}
		*reply = IncIDReply{
			ID: val,
		}
		return nil
	}

	s.mLock.Lock()
	defer s.mLock.Unlock()
	return s.emitCommand(vr.Command{
		Kind:  "IncID",
		Value: args,
	})
}

// ---
// internal write methods
// ---
func (s *Store) put(args PutArgs) {
	s.m[args.Key] = args.Val
}
func (s *Store) del(args DelArgs) {
	delete(s.m, args.Key)
}
func (s *Store) incid(args IncIDArgs) IncIDReply {
	s.incIDs[args.TableName]++
	return IncIDReply{
		ID: fmt.Sprintf("%d", s.incIDs[args.TableName]),
	}
}
func (s *Store) command(cmd vr.Command) {
	switch cmd.Kind {
	case "Put":
		{
			putArg, ok := cmd.Value.(PutArgs)
			if ok {
				s.put(putArg)
				return
			}
			break
		}
	case "Del":
		{
			arg, ok := cmd.Value.(DelArgs)
			if ok {
				s.del(arg)
				return
			}
			break
		}
	case "IncID":
		{
			arg, ok := cmd.Value.(IncIDArgs)
			if ok {
				s.incid(arg)
				return
			}
			break
		}
	}
}
func (s *Store) reset() {
	s.m = map[string]string{}
	s.incIDs = map[string]int{}
}

func (s *Store) processCommand(cmd vr.Command) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	s.command(cmd)
}

func (s *Store) replaceWithCommands(cmd []vr.Command) {
	s.mLock.Lock()
	defer s.mLock.Unlock()

	s.reset()
	for _, cmd := range cmd {
		s.command(cmd)
	}
}
