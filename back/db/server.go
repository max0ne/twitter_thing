package db

import (
	"net"
	"net/rpc"

	"github.com/max0ne/twitter_thing/back/config"
)

// Server serves a db store
type Server struct {
	store   *Store
	inbound *net.TCPListener
	port    string
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

	store := NewStore()
	if err := rpc.Register(store); err != nil {
		return nil, err
	}

	return &Server{
		store:   store,
		inbound: inbound,
		port:    config.DBPort,
	}, nil
}

// Start start hosting db
func (server *Server) Start() error {
	go rpc.Accept(server.inbound)
	return nil
}

// Port the port this server is listening on
func (server *Server) Port() string {
	return server.port
}
