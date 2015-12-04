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
    Title string `bson:"title" xml:"title"`
    Subtitle string `bson:"subtitle" xml:"subtitle"`
    Updated time.Time `bson:"updated" xml:"-"`
    URL string `bson:"url" xml:"-"`
    Icon string `bson:"icon" xml:"icon"`
    Twitter string `bson:"twitter" xml:"twitter"`
    Entries []Entry `bson:"entries" xml:"entry"`
    Name string `bson:"name" xml:"name"`

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