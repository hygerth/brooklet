package feedparser

import (
    "encoding/xml"
    "github.com/hygerth/brooklet/utils"
    "time"
)

type AtomFeed struct {
    XMLName  xml.Name  `xml:"feed"`
    Title    string    `xml:"title"`
    Subtitle string    `xml:"subtitle"`
    ID       string    `xml:"id"`
    Updated  time.Time `xml:"updated"`
    Link     []Link      `xml:"link"`
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

type RSSFeed struct {
    XMLName       xml.Name `xml:"rss"`
    Version       string   `xml:"version,attr"`
    Title         string   `xml:"channel>title"`
    Link          string   `xml:"channel>link"`
    Description   string   `xml:"channel>description"`
    LastBuildDate string   `xml:"channel>lastBuildDate"`
    Items         []Item   `xml:"channel>item"`
}

type Item struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    GUID        GUID   `xml:"guid"`
    Description string `xml:"description"`
    Content     string `xml:"encoded"`
    PubDate     string `xml:"pubDate"`
}

func ParseURL(url string) (AtomFeed, error) {
    var atomfeed AtomFeed
    page, err := utils.GetPage(url)
    if err != nil {
        return atomfeed, err
    }
    var rssfeed RSSFeed
    err = rssfeed.ParseFeed(page)
    if err == nil {
        af, _ := rssfeed.ConvertToAtomFeed()
        return af, nil
    }
    err = atomfeed.ParseFeed(page)
    if err != nil {
        return atomfeed, err
    }
    return atomfeed, err
}

func (atomfeed *AtomFeed) ParseFeed(feed []byte) error {
    err := xml.Unmarshal(feed, &atomfeed)
    return err
}

func (rssfeed *RSSFeed) ParseFeed(feed []byte) error {
    err := xml.Unmarshal(feed, &rssfeed)
    return err
}

func (rssfeed *RSSFeed) ConvertToAtomFeed() (AtomFeed, error) {
    var atomfeed AtomFeed
    t, _ := utils.ParseTimeString(rssfeed.LastBuildDate)
    atomfeed.Title = rssfeed.Title
    atomfeed.Link[0].Href = rssfeed.Link
    atomfeed.Link[0].Rel = "alternate"
    atomfeed.Subtitle = rssfeed.Description
    atomfeed.Updated = t
    var entries []Entry
    for _, item := range rssfeed.Items {
        var entry Entry
        it, _ := utils.ParseTimeString(item.PubDate)
        entry.Title = item.Title
        entry.Summary = item.Description
        entry.Content = item.Content
        entry.Updated = it
        entry.Link.Href = item.Link
        entry.ID = item.Link
        entry.GUID.Content = item.GUID.Content
        entry.GUID.IsPermaLink = item.GUID.IsPermaLink
        entries = append(entries, entry)
    }
    atomfeed.Entries = entries
    return atomfeed, nil
}
