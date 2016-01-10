package feedparser

import (
    "bytes"
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

type RSSFeed struct {
    XMLName       xml.Name  `xml:"rss"`
    Version       string    `xml:"version,attr"`
    Channel       RSSChannel `xml:"channel"`
}

type RSSChannel struct {
    Title         string   `xml:"title"`
    Link          string   `xml:"RssDefault link"`
    Description   string   `xml:"description"`
    LastBuildDate string   `xml:"lastBuildDate"`
    Items         []Item   `xml:"item"`
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
    d := xml.NewDecoder(bytes.NewReader(feed))
    d.DefaultSpace = "RssDefault"
    err := d.Decode(&rssfeed)
    return err
}

func (rssfeed *RSSFeed) ConvertToAtomFeed() (AtomFeed, error) {
    var atomfeed AtomFeed
    t, _ := utils.ParseTimeString(rssfeed.Channel.LastBuildDate)
    atomfeed.Title = rssfeed.Channel.Title
    link := Link{Href: rssfeed.Channel.Link, Rel: "alternate"}
    atomfeed.Link = append(atomfeed.Link, link)
    atomfeed.Subtitle = rssfeed.Channel.Description
    atomfeed.Updated = t
    var entries []Entry
    for _, item := range rssfeed.Channel.Items {
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
