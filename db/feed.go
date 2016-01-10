package db

import (
    "github.com/hygerth/brooklet/structure"
    "gopkg.in/mgo.v2/bson"
    "strings"
    "time"
)

func AddFeed(url string) (structure.Feed, error) {
    var feed structure.Feed
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return feed, err
    }
    c := session.DB(db).C("feeds")
    count, err := c.Find(bson.M{"url": url}).Count()
    if err != nil {
        return feed, err
    }
    if count > 0 {
        err = c.Find(bson.M{"url": url}).One(&feed)
        return feed, err
    }
    name := structure.GenerateFeedName(url)
    feed = structure.Feed{ID: bson.NewObjectId(), ChangeFrequency: 1, LastUpdate: time.Now().Add(-48 * time.Hour), URL: url, Name: name}
    err = c.Insert(feed)
    return feed, err
}

func RemoveFeed(url string) error {
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return err
    }
    c := session.DB(db).C("feeds")
    err = c.Remove(bson.M{"url": url})
    return err
}

func GetFeedByURL(url string) (structure.Feed, error) {
    var feed structure.Feed
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return feed, err
    }
    c := session.DB(db).C("feeds")
    err = c.Find(bson.M{"url": url}).One(&feed)
    feed = sortEntriesInFeed(feed)
    return feed, err
}

func GetFeedByName(name string) (structure.Feed, error) {
    var feed structure.Feed
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return feed, err
    }
    c := session.DB(db).C("feeds")
    err = c.Find(bson.M{"name": name}).One(&feed)
    feed = sortEntriesInFeed(feed)
    return feed, err
}

func UpdateFeedInformation(title, subtitle, icon, twitter, url string, updated time.Time) error {
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return err
    }
    c := session.DB(db).C("feeds")
    feed, err := GetFeedByURL(url)
    if err != nil {
        return err
    }
    feed.Title = title
    feed.Subtitle = subtitle
    feed.Icon = icon
    feed.Twitter = twitter
    feed.Updated = updated
    _, err = c.UpsertId(feed.ID, feed)
    return err
}

func UpdateFeedEntryList(entries []structure.Entry, url string) error {
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return err
    }
    c := session.DB(db).C("feeds")
    feed, err := GetFeedByURL(url)
    if err != nil {
        return err
    }
    for _, entry := range entries {
        exists := false
        for j, feedentry := range feed.Entries {
            if strings.Compare(feedentry.PermaLink, entry.PermaLink) == 0 {
                exists = true
                if !feedentry.Updated.Equal(entry.Updated) {
                    feed.Entries[j] = entry
                }
            }
        }
        if !exists {
            feed.Entries = append(feed.Entries, entry)
        }
    }
    _, err = c.UpsertId(feed.ID, feed)
    return err
}

func UpdateFeedChangeFrequency(url string, changefrequency float64) error {
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return err
    }
    c := session.DB(db).C("feeds")
    feed, err := GetFeedByURL(url)
    if err != nil {
        return err
    }
    feed.ChangeFrequency = changefrequency
    _, err = c.UpsertId(feed.ID, feed)
    return err
}

func UpdateFeedLastUpdate(url string, lastupdate time.Time) error {
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return err
    }
    c := session.DB(db).C("feeds")
    feed, err := GetFeedByURL(url)
    if err != nil {
        return err
    }
    feed.LastUpdate = lastupdate
    _, err = c.UpsertId(feed.ID, feed)
    return err
}

func GetAllFeeds() ([]structure.Feed, error) {
    var feeds []structure.Feed
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return feeds, err
    }
    c := session.DB(db).C("feeds")
    err = c.Find(bson.M{}).All(&feeds)
    feeds = sortEntriesInFeeds(feeds)
    return feeds, err
}

func GetFeedsInURLList(urls []string) ([]structure.Feed, error) {
    var feeds []structure.Feed
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return feeds, err
    }
    c := session.DB(db).C("feeds")
    err = c.Find(bson.M{"url": bson.M{"$in": urls}}).All(&feeds)
    feeds = sortEntriesInFeeds(feeds)
    return feeds, err
}

func GetEntryByArticleID(id string) (structure.Entry, error) {
    var entry structure.Entry
    var feed structure.Feed
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return entry, err
    }
    c := session.DB(db).C("feeds")
    err = c.Find(bson.M{"entries.articleid": id}).One(&feed)
    if err != nil {
        return entry, err
    }
    for _, e := range feed.Entries {
        if e.ArticleID == id {
            entry = e
            break
        }
    }
    return entry, err
}

func HasEntryByArticleID(id string) (bool, error) {
    session, err := Connect()
    defer session.Close()
    if err != nil {
        return false, err
    }
    c := session.DB(db).C("feeds")
    count, err := c.Find(bson.M{"entries.articleid": id}).Count()
    if err != nil {
        return false, err
    }
    if count > 0 {
        return true, nil
    }
    return false, nil
}

func sortEntriesInFeeds(feeds []structure.Feed) []structure.Feed {
    for i, feed := range feeds {
        feeds[i] = sortEntriesInFeed(feed)
    }
    return feeds
}

func sortEntriesInFeed(feed structure.Feed) structure.Feed {
    feed.Entries = structure.SortEntriesByDate(feed.Entries)
    return feed
}
