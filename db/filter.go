package db

import (
    "github.com/hygerth/brooklet/structure"
    "gopkg.in/mgo.v2/bson"
)

func AddFilter(filter string) (structure.Filter, error) {
    var newfilter structure.Filter
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return newfilter, err
    }
    c := session.DB(db).C("filter")
    count, err := c.Find(bson.M{"filter": filter}).Count()
    if err != nil {
        return newfilter, err
    }
    if count > 0 {
        err = c.Find(bson.M{"filter": filter}).One(&filter)
        return newfilter, err
    }
    newfilter = structure.Filter{ID: bson.NewObjectId(), Filter: filter}
    err = c.Insert(newfilter)
    return newfilter, err
}

func RemoveFilter(filter string) error {
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return err
    }
    c := session.DB(db).C("filter")
    err = c.Remove(bson.M{"filter": filter})
    return err
}

func GetFilter() ([]structure.Filter, error) {
    var filter []structure.Filter
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return filter, err
    }
    c := session.DB(db).C("filter")
    err = c.Find(bson.M{}).All(&filter)
    return filter, err
}