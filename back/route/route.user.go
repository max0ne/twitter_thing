package model

import (
	"encoding/json"

	"github.com/max0ne/twitter_thing/back/db"
)

// GetUser - -
func GetUser(tid string, table *db.Table) (*User, error) {
	var User User
	if err := json.Unmarshal([]byte(table.Get(tid)), &User); err != nil {
		return nil, err
	}
	return &User, nil
}
