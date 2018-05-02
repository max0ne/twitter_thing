package util

import (
	"encoding/json"
	"log"
	"net"
	"net/rpc"
	"os"
)

// Contains - -
func Contains(vals []string, aVal string) bool {
	for _, v := range vals {
		if v == aVal {
			return true
		}
	}
	return false
}

// Remove - -
func Remove(vals []string, aVal string) []string {
	newVals := []string{}
	for _, v := range vals {
		if v == aVal {
			continue
		}
		newVals = append(newVals, v)
	}
	return newVals
}

// JSONMarshel - -
func JSONMarshel(val interface{}) string {
	str, _ := json.Marshal(val)
	return string(str)
}

// GetEnvMust get env, crashes if env key not set
func GetEnvMust(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal("env key ", key, " missing")
	}
	return val
}

// NewRPC make an rpc server
func NewRPC(url string, rcvr interface{}) (*net.TCPListener, *rpc.Server, error) {
	addy, err := net.ResolveTCPAddr("tcp", url)
	if err != nil {
		return nil, nil, err
	}
	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		return nil, nil, err
	}
	// rcp server
	rpcServer := rpc.NewServer()
	if err := rpcServer.Register(rcvr); err != nil {
		return nil, nil, err
	}
	return inbound, rpcServer, nil
}
