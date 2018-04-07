package model

import (
	"encoding/json"
	"strings"

	"github.com/max0ne/twitter_thing/back/db"
)

// User - -
type User struct {
	Uname    string `json:"uname"`
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
func GetUser(uname string, table *db.Table) (*User, error) {
	if !table.Has(uname) {
		return nil, nil
	}

	var User User
	if err := json.Unmarshal([]byte(table.Get(uname)), &User); err != nil {
		return nil, err
	}
	return &User, nil
}

// SaveUser - -
func SaveUser(user User, table *db.Table) error {
	bytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	table.Put(user.Uname, string(bytes))
	return nil
}

// DeleteUser - -
func DeleteUser(user User, table *db.Table) error {
	table.Del(user.Uname)
	return nil
}

// Follow - -
func Follow(user User, vid string, followTable *db.Table) error {
	followingString := followTable.Get(user.Uname)
	following := strings.Split(followingString, ",")
	for _, id := range following {
		if id == vid {
			return nil
		}
	}

	followTable.Put(user.Uname, followingString+","+vid)
	return nil
}

// UnfollowUser - -
func UnfollowUser(user User, vid string, followTable *db.Table) error {
	followingString := followTable.Get(user.Uname)
	following := strings.Split(followingString, ",")
	newFollowing := []string{}
	for _, id := range following {
		if id != vid {
			newFollowing = append(newFollowing, id)
		}
	}

	followTable.Put(user.Uname, strings.Join(newFollowing, ","))
	return nil
}
