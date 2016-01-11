package structure

import (
    "sort"
    "time"
)

type Entry struct {
    Title string `bson:"title" xml:"title,omitempty"`
    Summary string `bson:"summary" xml:"summary,omitempty"`
    Content string `bson:"content" xml:"-"`
    ID string `bson:"id" xml:"-"`
    Updated time.Time `bson:"updated" xml:"published,omitempty"`
    PermaLink string `bson:"permalink" xml:"url,omitempty"`
    Author string `bson:"author" xml:"author,omitempty"`
    Twitter string `bson:"twitter" xml:"twitter,omitempty"`
    HasImage bool `bson:"hasimage" xml:"hasimage,omitempty"`
    ImageRotation string `bson:"imagerotation" xml:"imagerotation,omitempty"`
    ArticleID string `bson:"articleid" xml:"id,omitempty"`
}

type Entries []Entry

func (e Entries) Len() int {
    return len(e)
}

func (e Entries) Less(i, j int) bool {
    return e[i].Updated.After(e[j].Updated)
}

func (e Entries) Swap(i, j int) {
    e[i], e[j] = e[j], e[i]
}

// SortEntriesByDate sorts the entries in a list by the date of which they
// were updated on
func SortEntriesByDate(entries []Entry) []Entry {
    entriessorted := make(Entries, 0, len(entries))
    for _, entry := range entries {
        entriessorted = append(entriessorted, entry)
    }
    sort.Sort(entriessorted)
    var en []Entry
    for _, entry := range entriessorted {
        en = append(en, entry)
    }
    return en
}