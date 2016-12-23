package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var globalUserStore UserStore

//var globalUsernameStore []string
//var globalUserEmailStore []string

type UserStore interface {
	Find(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindByUsername(string) (*User, error)
	Save(User) error
}

func InitUserStore() {
	store, err := NewFileUserStore("./data/users.json")
	if err != nil {
		panic(fmt.Errorf("Error creating user store: %s", err))
	}
	globalUserStore = store
}

type FileUserStore struct {
	filename string
	Users    map[string]User
}

func (store FileUserStore) Save(user User) error {
	store.Users[user.ID] = user
	//globalUsernameStore = append(globalUsernameStore, user.UserName)
	//globalUserEmailStore = append(globalUserEmailStore, user.Email)

	contents, err := json.MarshalIndent(store, "", "   ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(store.filename, contents, 0660)
	if err != nil {
		return err
	}
	return nil
}

func (store FileUserStore) Find(id string) (*User, error) {
	user, ok := store.Users[id]
	if ok {
		return &user, nil
	}
	return nil, nil
}

func (store FileUserStore) FindByUsername(username string) (*User, error) {
	if username == "" {
		return nil, nil
	}
	for _, user := range store.Users {
		if strings.ToLower(username) == strings.ToLower(user.UserName) {
			return &user, nil
		}
	}
	/*for _, usernameGlobal := range globalUsernameStore {
		if strings.ToLower(username) == strings.ToLower(usernameGlobal) {
			return &User{}, nil
		}
	}*/
	return nil, nil
}

func (store FileUserStore) FindByEmail(email string) (*User, error) {
	if email == "" {
		return nil, nil
	}

	for _, user := range store.Users {
		if strings.ToLower(email) == strings.ToLower(user.Email) {
			return &user, nil
		}
	}
	/*for _, userEmailGlobal := range globalUserEmailStore {
		if strings.ToLower(email) == strings.ToLower(userEmailGlobal) {
			return &User{}, nil
		}
	}*/
	return nil, nil
}

func NewFileUserStore(filename string) (*FileUserStore, error) {
	store := &FileUserStore{
		Users:    map[string]User{},
		filename: filename,
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		// If it's a matter of the file not existing, that's ok
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}
	err = json.Unmarshal(contents, store)
	if err != nil {
		return nil, err
	}
	/*for _, usr := range store.Users {
		globalUsernameStore = append(globalUsernameStore, usr.UserName)
		globalUserEmailStore = append(globalUserEmailStore, usr.Email)
	}*/
	return store, nil
}
