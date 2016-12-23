package main

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	dbName         = "gophr"
	collectionName = "images"
	pageSize       = 25
)

type DBImageStore struct {
	Session *mgo.Session
}

type ImageStore interface {
	Save(image *Image) error
	Find(id string) (*Image, error)
	FindAll(offset int) ([]Image, error)
	FindAllByUser(user *User, offset int) ([]Image, error)
}

func NewDBImageStore() *DBImageStore {
	return &DBImageStore{
		Session: mongoSession.Copy(),
	}
}

func (store *DBImageStore) Save(image *Image) error {
	err := store.Session.DB(dbName).C(collectionName).Insert(&image)
	if err != nil {
		return err
	}
	return nil
}

func (store *DBImageStore) Find(id string) (*Image, error) {
	image := &Image{}
	err := store.Session.DB(dbName).C(collectionName).Find(bson.M{"_id": id}).One(image)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (store *DBImageStore) FindAll(offset int) ([]Image, error) {
	var results []Image
	err := store.Session.DB(dbName).C(collectionName).Find(nil).Skip(offset).Limit(pageSize).Sort("-created_at").All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (store *DBImageStore) FindAllByUser(user *User, offset int) ([]Image, error) {
	var results []Image
	err := store.Session.DB(dbName).C(collectionName).Find(bson.M{"user_id": user.ID}).Skip(offset).Limit(pageSize).Sort("-created_at").All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (store *DBImageStore) Close() {
	store.Session.Close()
}
