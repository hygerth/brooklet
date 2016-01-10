package atom

import (
    "encoding/xml"
    "time"
)

type AtomFeed struct {
    XMLName  xml.Name  `xml:"feed"`
    XMLNS    string    `xml:"xmlns,attr"`
    Title    string    `xml:"title"`
    Subtitle string    `xml:"subtitle"`
    ID       string    `xml:"id"`
    Updated  time.Time `xml:"updated"`
    Link     []Link    `xml:"link"`
    Author   Author    `xml:"author"`
    Entries  []Entry   `xml:"entry"`
}

type Link struct {
    Href string `xml:"href,attr"`
    Rel string `xml:"rel,attr"`
}

type Author struct {
    Name  string `xml:"name"`
    Email string `xml:"email"`
}

type Entry struct {
    Title   string    `xml:"title"`
    Summary string    `xml:"summary"`
    Content string    `xml:"content"`
    ID      string    `xml:"id"`
    GUID    GUID      `xml:"guid"`
    Updated time.Time `xml:"updated"`
    Link    Link      `xml:"link"`
    Author  Author    `xml:"author"`
}

type GUID struct {
    Content string `xml:",chardata"`
    IsPermaLink bool `xml:"isPermaLink,attr"`
}

func CreateAtomFeed(title, subtitle, id, authorname, authoremail string, updated time.Time) AtomFeed {
    var atomfeed AtomFeed
    atomfeed.XMLNS = "http://www.w3.org/2005/Atom"
    atomfeed.Title = title
    atomfeed.Subtitle = subtitle
    atomfeed.ID = id
    atomfeed.Updated = updated
    atomfeed.Author = Author{Name: authorname, Email: authoremail}
    return atomfeed
}

func (a *AtomFeed) AddEntry(title, summary, content, id, guidcontent, link, authorname, authoremail string, updated time.Time) {
    var entry Entry
    entry.Title = title
    entry.Summary = summary
    entry.Content = content
    entry.ID = id
    entry.GUID = GUID{Content: guidcontent, IsPermaLink: false}
    entry.Updated = updated
    entry.Link = Link{Href: link, Rel: "alternate"}
    entry.Author = Author{Name: authorname, Email: authoremail}
    a.Entries = append(a.Entries, entry)
}