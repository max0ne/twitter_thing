package db

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/max0ne/twitter_thing/back/config"
)

// Server serves a db store
type Server struct {
	store     *Store
	inbound   *net.TCPListener
	rpcServer *rpc.Server
	port      string
}

// NewServer make a db hosting server
func NewServer(config config.Config) (*Server, error) {
	addy, err := net.ResolveTCPAddr("tcp", config.DBURL())
	if err != nil {
		return nil, err
	}

	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		return nil, err
	}

	// storage
	store := NewStore()

	// rcp server
	rpcServer := rpc.NewServer()
	if err := rpcServer.Register(store); err != nil {
		return nil, err
	}

	return &Server{
		store:     store,
		inbound:   inbound,
		rpcServer: rpcServer,
		port:      config.DBPort,
	}, nil
}

// Start start hosting db
func (server *Server) Start() error {
	fmt.Println("db rpc server hosting on", server.Port())
	go server.rpcServer.Accept(server.inbound)
	return nil
}

// Port the port this server is listening on
func (server *Server) Port() string {
	return server.port
}
