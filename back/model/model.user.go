package model

import (
	"encoding/json"

	"github.com/max0ne/twitter_thing/back/db"
)

// User - -
type User struct {
	Uname    string `json:"uname"`
	UID      string `json:"uid"`
	Password string `json:"password"`
}

// NewUser - -
func NewUser(uname, password string) User {
	return User{
		Uname:    uname,
		Password: password,
	}
}

// GetUser - -
func GetUser(tid string, table *db.Table) (*User, error) {
	var User User
	if err := json.Unmarshal([]byte(table.Get(tid)), &User); err != nil {
		return nil, err
	}
	return &User, nil
}
