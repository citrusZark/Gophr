package main

import mgo "gopkg.in/mgo.v2"

const (
	host = "localhost"
)

var mongoSession *mgo.Session

func InitMongoDB() {
	db, err := NewMongoDBSession(host)
	if err != nil {
		panic(err)
	}
	mongoSession = db
}

func NewMongoDBSession(con string) (*mgo.Session, error) {
	db, err := mgo.Dial(con)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CopyMongoDBSession(db *mgo.Session) *mgo.Session {
	dbClone := db.Copy()
	return dbClone
}

func CloseMongoDBSession(db *mgo.Session) {
	db.Close()
}
