package db

import (
	"encoding/json"
	"net/rpc"

	"github.com/max0ne/twitter_thing/back/config"
)

// Client db client
type Client struct {
	rpc *rpc.Client
}

// Table - -
type Table struct {
	Name     string
	dbClient *Client
}

// NewClient build a db client
func NewClient(config config.Config) (*Client, error) {
	client, err := rpc.Dial("tcp", config.DBURL())
	if err != nil {
		return nil, err
	}
	return &Client{
		rpc: client,
	}, nil
}

// NewTable new a table
func (client *Client) NewTable(name string) *Table {
	return &Table{
		Name:     name,
		dbClient: client,
	}
}

// GetM - for debug only
func (client *Client) GetM() (map[string]string, error) {
	reply := map[string]string{}
	err := client.rpc.Call("Store.GetM", nil, &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// Get - -
func (client *Client) Get(key string) (string, error) {
	reply := GetReply{}
	err := client.rpc.Call("Store.Get", GetArgs{Key: key}, &reply)
	if err != nil {
		return "", err
	}
	return reply.Val, nil
}

// Put - -
func (client *Client) Put(key, val string) error {
	return client.rpc.Call("Store.Put", PutArgs{Key: key, Val: val}, nil)
}

// Del - -
func (client *Client) Del(key string) error {
	return client.rpc.Call("Store.Del", DelArgs{Key: key}, nil)
}

// Has - -
func (client *Client) Has(key string) (bool, error) {
	reply := HasReply{}
	err := client.rpc.Call("Store.Has", GetArgs{Key: key}, &reply)
	if err != nil {
		return false, err
	}
	return reply.Has, nil
}

// IncID - -
func (client *Client) IncID(tableName string) (string, error) {
	reply := IncIDReply{}
	err := client.rpc.Call("Store.IncID", IncIDArgs{TableName: tableName}, &reply)
	if err != nil {
		return "", err
	}
	return reply.ID, nil
}

// Put - -
func (t *Table) Put(key, val string) error {
	return t.dbClient.Put(t.key(key), val)
}

// PutObj - -
func (t *Table) PutObj(key string, val interface{}) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return t.Put(key, string(bytes))
}

// Del - -
func (t *Table) Del(key string) error {
	return t.dbClient.Del(t.key(key))
}

// Get - -
func (t *Table) Get(key string) (string, error) {
	return t.dbClient.Get(t.key(key))
}

// GetObj - -
func (t *Table) GetObj(key string, target interface{}) error {
	jsonString, err := t.Get(key)
	if err != nil {
		return err
	}
	if jsonString == "" {
		return nil
	}
	return json.Unmarshal([]byte(jsonString), target)
}

// Has has key
func (t *Table) Has(key string) (bool, error) {
	return t.dbClient.Has(t.key(key))
}

func (t *Table) key(key string) string {
	return t.Name + "_" + key
}

// IncID get auto increment id
func (t *Table) IncID() (string, error) {
	return t.dbClient.IncID(t.Name)
}
