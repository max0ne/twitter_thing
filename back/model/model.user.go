package model

import (
	"strings"

	"github.com/max0ne/twitter_thing/back/db"
	"github.com/max0ne/twitter_thing/back/util"
)

const keyNewRegisteredUsers = "keyNewRegisteredUsers"

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

// GetUser - remove password field
func GetUser(uname string, table *db.Table) (*User, error) {
	user, err := GetUserWithPassword(uname, table)
	if user != nil {
		user.Password = ""
	}
	return user, err
}

// GetUserWithPassword - -
func GetUserWithPassword(uname string, table *db.Table) (*User, error) {
	hasUser, err := table.Has(uname)
	if err != nil {
		return nil, err
	}
	if !hasUser {
		return nil, nil
	}

	var User User
	err = table.GetObj(uname, &User)
	if err != nil {
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
func SaveUser(user User, userTable *db.Table, miscTable *db.Table) error {
	if err := userTable.PutObj(user.Uname, user); err != nil {
		return err
	}
	return PutNewRegisteredUserName(user.Uname, miscTable)
}

// DeleteUser - -
func DeleteUser(user User, table *db.Table) error {
	table.Del(user.Uname)
	return nil
}

// Follow - -
func Follow(user User, vname string, followingTable, followerTable *db.Table) error {

	insert := func(key, val string, table *db.Table) error {
		vstring, err := table.Get(key)
		if err != nil {
			return err
		}
		if util.Contains(strings.Split(vstring, ","), val) {
			return nil
		}
		return table.Put(key, vstring+","+val)
	}

	if err := insert(user.Uname, vname, followingTable); err != nil {
		return err
	}
	return insert(vname, user.Uname, followerTable)
}

// Unfollow - -
func Unfollow(user User, vname string, followingTable, followerTable *db.Table) error {
	remove := func(key, val string, table *db.Table) error {
		ss, err := table.Get(key)
		if err != nil {
			return err
		}
		return table.Put(key, strings.Join(util.Remove(strings.Split(ss, ","), val), ","))
	}
	if err := remove(user.Uname, vname, followingTable); err != nil {
		return err
	}
	return remove(vname, user.Uname, followerTable)
}

// GetFollowers - -
func GetFollowers(vname string, followerTable *db.Table) ([]string, error) {
	ss, err := followerTable.Get(vname)
	if err != nil {
		return nil, err
	}
	return strings.Split(ss, ","), nil
}

// GetFollowing - -
func GetFollowing(vname string, followerTable *db.Table) ([]string, error) {
	return GetFollowers(vname, followerTable)
}

// GetNewRegisteredUserNames get list of recently registered users' names
func GetNewRegisteredUserNames(miscTable *db.Table) ([]string, error) {
	newRegisteredUserNames := []string{}
	err := miscTable.GetObj(keyNewRegisteredUsers, &newRegisteredUserNames)
	if err != nil {
		return nil, err
	}
	return newRegisteredUserNames, nil
}

// PutNewRegisteredUserName append a new value to newly registered users list
func PutNewRegisteredUserName(uname string, miscTable *db.Table) error {
	has, err := miscTable.Has(keyNewRegisteredUsers)
	if err != nil {
		return err
	}

	newRegisteredUserNames := []string{}
	if has {
		err = miscTable.GetObj(keyNewRegisteredUsers, &newRegisteredUserNames)
		if err != nil {
			return err
		}
	}
	newRegisteredUserNames = append(newRegisteredUserNames, uname)
	// only keep recent 10 items
	if len(newRegisteredUserNames) > 10 {
		newRegisteredUserNames = newRegisteredUserNames[len(newRegisteredUserNames)-10:]
	}
	return miscTable.PutObj(keyNewRegisteredUsers, newRegisteredUserNames)
}
