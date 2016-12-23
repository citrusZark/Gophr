package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var globalSessionStore SessionStore

func InitSessionStore() {
	sessionStore, err := NewFileSessionStore("./data/sessions.json")
	if err != nil {
		panic(fmt.Errorf("Error creating session store: %s", err))
	}
	globalSessionStore = sessionStore
}

type SessionStore interface {
	Find(string) (*Session, error)
	Save(Session) error
	Delete(*Session) error
}

type FileSessionStore struct {
	filename string
	Sessions map[string]Session
}

func (store FileSessionStore) Save(session Session) error {
	store.Sessions[session.ID] = session

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

func (store FileSessionStore) Find(id string) (*Session, error) {
	session, ok := store.Sessions[id]
	if ok {
		return &session, nil
	}
	return nil, nil
}

func (store FileSessionStore) Delete(sesion *Session) error {
	delete(store.Sessions, sesion.ID)
	contents, err := json.MarshalIndent(store, "", "   ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(store.filename, contents, 0660)
}

func NewFileSessionStore(name string) (*FileSessionStore, error) {
	store := &FileSessionStore{
		Sessions: map[string]Session{},
		filename: name,
	}

	contents, err := ioutil.ReadFile(name)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, nil
	}
	err = json.Unmarshal(contents, store)
	if err != nil {
		return nil, err
	}
	return store, nil
}
