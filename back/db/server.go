package db

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/max0ne/twitter_thing/back/config"
	"github.com/max0ne/twitter_thing/back/util"
	"github.com/max0ne/twitter_thing/back/vr"
)

// Server serves a db store
type Server struct {
	store    *Store
	vrServer *vr.PBServer
	port     string
}

func (server *Server) emitVRCommand(cmd interface{}) {
	vr.PushCommand(server.vrServer, cmd)
}

// RunServer make a db hosting server
func RunServer(config config.Config) (*Server, error) {

	server := Server{}

	// 1. setup db
	store := NewStore(server.emitVRCommand, config.VRPrimary == config.VRMe())
	dbConn, dbRPCServer, err := util.NewRPC(config.DBURL(), store)
	if err != nil {
		return nil, err
	}

	// 2. setup vr model
	vrServer := vr.Make(store.processWriteCommand)

	// 3. setup vr rpc
	vrConn, vrRPCServer, err := util.NewRPC(config.VRURL(), vrServer)
	if err != nil {
		return nil, err
	}

	// 4. start hosting rpc servers
	fmt.Println("vr rpc server listening on", config.VRURL())
	go vrRPCServer.Accept(vrConn)

	// 5. sleep for a while to wait for all vr server to start listening
	<-time.After(time.Second)

	// 6. connect all vr peers
	rpcClients, err := connectVRPeers(config)
	if err != nil {
		return nil, err
	}

	// 7. start vr logic
	vr.Start(vrServer, rpcClients, config.VRMe())

	// 8. start hosting db server
	fmt.Println("db rpc server listening on", config.DBURL())
	go dbRPCServer.Accept(dbConn)

	server.store = store
	server.vrServer = vrServer
	server.port = config.DBPort
	return &server, nil
}

func connectVRPeers(config config.Config) ([]*rpc.Client, error) {
	clients := []*rpc.Client{}
	for _, peerURL := range config.VRPeerURLs {
		// don't connect me, add a dummy client just to align array indices
		if peerURL == config.VRURL() {
			clients = append(clients, &rpc.Client{})
			continue
		}
		// connect peer
		client, err := rpc.Dial("tcp", peerURL)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

// Port the port this server is listening on
func (server *Server) Port() string {
	return server.port
}
