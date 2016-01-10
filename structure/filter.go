package structure

import (
    "gopkg.in/mgo.v2/bson"
    "strings"
)

type Filter struct {
    ID bson.ObjectId `bson:"_id,omitempty" xml:"-"`
    Filter string `bson:"filter" xml:"filter"`
}

func FilterEntries(entries []Entry, filter []Filter) []Entry {
    if len(filter) == 0 {
        return entries
    }
    var filteredentries []Entry
    for _, entry := range entries {
        clean := true
        for _, f := range filter {
            if strings.Contains(entry.Title, f.Filter) || strings.Contains(entry.Summary, f.Filter) || strings.Contains(entry.Content, f.Filter) {
                clean = false
                break
            }
        }
        if clean {
            filteredentries = append(filteredentries, entry)
        }
    }
    return filteredentries
}