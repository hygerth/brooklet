package skywalker

import (
    "github.com/hygerth/brooklet/db"
    "github.com/hygerth/brooklet/feedparser"
    "github.com/hygerth/brooklet/siteparser"
    "github.com/hygerth/brooklet/structure"
    "github.com/hygerth/brooklet/utils"
    "strings"
    "sync"
    "time"
)

func processFeeds(feeds []structure.Feed) {
    in := genInChan(feeds)
    var wg sync.WaitGroup
    for i := 0; i < 20; i++ {
        wg.Add(1)
        updateWorker(in, &wg)
    }
    wg.Wait()
}

func genInChan(feeds []structure.Feed) <-chan structure.Feed {
    out := make(chan structure.Feed)
    go func() {
        for _, feed := range feeds {
            out <- feed
        }
        close(out)
    }()
    return out
}

func updateWorker(in <-chan structure.Feed, wg *sync.WaitGroup) {
    go func() {
        defer wg.Done()
        for feed := range in {
            UpdateFeed(feed)
        }
    }()
}

func UpdateFeed(feed structure.Feed) {
    update, err := feedparser.ParseURL(feed.URL)
    utils.Checkerr(err)
    var link string
    for _, li := range update.Link {
        if li.Rel == "alternate" {
            link = li.Href
        }
    }
    if len(link) == 0 {
        link = update.Link[0].Href
    }
    meta, err := siteparser.GetMetaForSite(link)
    utils.Checkerr(err)
    err = db.UpdateFeedInformation(update.Title, update.Subtitle, meta.Icon, meta.TwitterSite, feed.URL, update.Updated)
    utils.Checkerr(err)
    var entries []structure.Entry
    for _, entry := range update.Entries {
        exists := false
        for _, oldentry := range feed.Entries {
            if strings.Compare(entry.ID, oldentry.ID) == 0 && entry.Updated.Equal(oldentry.Updated) {
                exists = true
            }
        }
        if !exists {
            newentry := convertAtomEntryToDBEntry(entry)
            entries = append(entries, newentry)
        }
    }
    err = db.UpdateFeedEntryList(entries, feed.URL)
    utils.Checkerr(err)
    feed, err = db.GetFeedByURL(feed.URL)
    utils.Checkerr(err)
    sortedentries := structure.SortEntriesByDate(feed.Entries)
    var occasions []time.Time
    for _, entry := range sortedentries {
        occasions = append(occasions, entry.Updated)
    }
    changefrequency := calculateChangeFrequency(occasions)
    err = db.UpdateFeedChangeFrequency(feed.URL, changefrequency)
    utils.Checkerr(err)
    err = db.UpdateFeedLastUpdate(feed.URL, time.Now())
    utils.Checkerr(err)
}

func convertAtomEntryToDBEntry(entry feedparser.Entry) structure.Entry {
    var newentry structure.Entry
    newentry.Title = entry.Title
    newentry.Summary = entry.Summary
    if len(entry.Summary) >= len(entry.Content) {
        article, err := siteparser.GetArticleForSite(entry.Link.Href)
        utils.Checkerr(err)
        newentry.Content = article
    } else {
        newentry.Content = entry.Content
    }
    newentry.ID = entry.ID
    newentry.Updated = entry.Updated
    if !entry.GUID.IsPermaLink && len(entry.GUID.Content) > 0 {
        newentry.PermaLink = entry.GUID.Content
    } else {
        newentry.PermaLink = entry.Link.Href
    }
    newentry.Author = entry.Author.Name
    meta, err := siteparser.GetMetaForSite(newentry.PermaLink)
    utils.Checkerr(err)
    newentry.Twitter = meta.TwitterCreator
    filename, isPortrait, err := SyncImage(meta.Image)
    utils.Checkerr(err)
    newentry.Image.BaseFilename = filename
    switch isPortrait {
    case false: newentry.Image.Rotation = "landscape"
    default: newentry.Image.Rotation = "portrait"
    }
    return newentry
}

func calculateChangeFrequency(occasions []time.Time) float64 {
    var total float64
    for i := 0; i < len(occasions) - 1; i++ {
        total += occasions[i].Sub(occasions[i + 1]).Hours()
    }
    avg := total/float64(len(occasions))
    timesbetweenchangesperday := avg / 24
    return 1/timesbetweenchangesperday
}