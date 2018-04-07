package model

import (
	"encoding/json"
	"strings"

	"github.com/max0ne/twitter_thing/back/db"
	"github.com/max0ne/twitter_thing/back/util"
)

// User - -
type User struct {
	Uname    string `json:"username"`
	Password string `json:"password"`
}

// NewUser - -
func NewUser(uname, password string) User {
	return User{
		Uname:    uname,
		Password: password,
	}
}

// GetUser - remove password field
func GetUser(uname string, table *db.Table) (*User, error) {
	if !table.Has(uname) {
		return nil, nil
	}

	var User User
	if err := json.Unmarshal([]byte(table.Get(uname)), &User); err != nil {
		return nil, err
	}
	User.Password = ""
	return &User, nil
}

// GetUserWithPassword - -
func GetUserWithPassword(uname string, table *db.Table) (*User, error) {
	if !table.Has(uname) {
		return nil, nil
	}

	var User User
	if err := json.Unmarshal([]byte(table.Get(uname)), &User); err != nil {
		return nil, err
	}
	return &User, nil
}

// GetUsers - -
func GetUsers(unames []string, userTable *db.Table) []User {
	users := []User{}
	for _, uname := range unames {
		user, err := GetUser(uname, userTable)
		if user != nil && err == nil {
			users = append(users, *user)
		}
	}
	return users
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
func Follow(user User, vname string, followingTable, followerTable *db.Table) error {

	insert := func(key, val string, table *db.Table) {
		vstring := table.Get(key)
		if util.Contains(strings.Split(vstring, ","), val) {
			return
		}
		table.Put(key, vstring+","+val)
	}

	insert(user.Uname, vname, followingTable)
	insert(vname, user.Uname, followerTable)
	return nil
}

// Unfollow - -
func Unfollow(user User, vname string, followingTable, followerTable *db.Table) error {
	remove := func(key, val string, table *db.Table) {
		table.Put(key, strings.Join(util.Remove(strings.Split(table.Get(key), ","), val), ","))
	}
	remove(user.Uname, vname, followingTable)
	remove(vname, user.Uname, followerTable)
	return nil
}

// GetFollowers - -
func GetFollowers(vname string, followerTable *db.Table) []string {
	return strings.Split(followerTable.Get(vname), ",")
}

// GetFollowing - -
func GetFollowing(vname string, followerTable *db.Table) []string {
	return GetFollowers(vname, followerTable)
}
