package structure

import (
    "gopkg.in/mgo.v2/bson"
    "net/url"
    "regexp"
    "strings"
    "time"
)

var regex = regexp.MustCompile(`\W`)

type Feed struct {
    Title string `bson:"title" xml:"title,omitempty"`
    Subtitle string `bson:"subtitle" xml:"subtitle,omitempty"`
    Updated time.Time `bson:"updated" xml:"-"`
    URL string `bson:"url" xml:"-"`
    Icon string `bson:"icon" xml:"icon,omitempty"`
    Twitter string `bson:"twitter" xml:"twitter,omitempty"`
    Entries []Entry `bson:"entries" xml:"entry,omitempty"`
    Name string `bson:"name" xml:"name,omitempty"`

    ID bson.ObjectId `bson:"_id,omitempty" xml:"-"`
    ChangeFrequency float64 `bson:"changefrequency" xml:"-"`
    LastUpdate time.Time `bson:"lastupdate" xml:"-"`
}

func GenerateFeedName(feedurl string) string {
    u, _ := url.Parse(feedurl)
    s := u.Host + u.Path
    s = regex.ReplaceAllString(s, "")
    s = strings.ToLower(s)
    return s
}

func ExtractEntriesFromFeeds(feeds ...Feed) []Entry {
    var entries []Entry
    for _, feed := range feeds {
        entries = append(entries, feed.Entries...)
    }
    entries = SortEntriesByDate(entries)
    return entries
}
