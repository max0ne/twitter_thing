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
	store         *Store
	vrServer      *vr.PBServer
	dbPeerClients []*Client
	port          string
}

func (server *Server) emitVRCommand(cmd vr.Command) error {
	go server.vrServer.PushCommand(cmd, nil)
	return nil
}

// RunServer make a db hosting server
func RunServer(config config.Config) (*Server, error) {
	server := Server{}

	// 1. setup db
	store := NewStore(server.emitVRCommand, &server)
	dbConn, dbRPCServer, err := util.NewRPC(config.DBURL(), store)
	if err != nil {
		return nil, err
	}

	// 2. setup vr model
	vrServer := vr.Make(config.VRMe(), func(cmd vr.Command) {
		util.LogColor(config.VRMe())(config.VRMe(), "process command", cmd)
		store.processCommand(cmd)
	}, func(cmds []vr.Command) {
		util.LogColor(config.VRMe())(config.VRMe(), "replace commands")
		store.replaceWithCommands(cmds)
	})

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

	// 6.5. connect all db peers
	dbPeerClients, err := connectDBPeers(config)
	if err != nil {
		return nil, err
	}

	// 7. start vr logic
	vr.Start(vrServer, rpcClients)

	// 8. start hosting db server
	fmt.Println("db rpc server listening on", config.DBURL())
	go dbRPCServer.Accept(dbConn)

	server.store = store
	server.vrServer = vrServer
	server.port = config.DBPort
	server.dbPeerClients = dbPeerClients
	return &server, nil
}

func connectVRPeers(config config.Config) ([]*rpc.Client, error) {
	clients := []*rpc.Client{}
	for _, peerURL := range config.VRPeerURLs {
		// connect peer
		client, err := rpc.Dial("tcp", peerURL)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func connectDBPeers(config config.Config) ([]*Client, error) {
	clients := []*Client{}
	for _, peerURL := range config.DBPeerURLs {
		client, err := NewClient(peerURL)
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
