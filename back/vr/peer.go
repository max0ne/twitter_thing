package vr

import (
	"net/rpc"
)

type peer struct {
	rpcClient *rpc.Client
}

func (peer *peer) Call(serviceMethod string, args interface{}, reply interface{}) bool {
	err := peer.rpcClient.Call(serviceMethod, args, reply)
	return err == nil
}

func makePeers(rpcClients []*rpc.Client) []*peer {
	peers := []*peer{}
	for _, rpcClient := range rpcClients {
		peers = append(peers, &peer{rpcClient: rpcClient})
	}
	return peers
}
