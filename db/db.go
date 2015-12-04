package db

import (
    "gopkg.in/mgo.v2"
)

const (
    uri string = "mongodb://127.0.0.1:27017/brooklet"
    db  string = "brooklet"
)

// Connect establishes an connection to the database
func Connect() (*mgo.Session, error) {
    session, err := mgo.Dial(uri)
    session.SetSafe(&mgo.Safe{})
    return session, err
}