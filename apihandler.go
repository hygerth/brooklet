package brooklet

import (
    "github.com/gorilla/mux"
    "github.com/hygerth/brooklet/atom"
    "github.com/hygerth/brooklet/db"
    "github.com/hygerth/brooklet/structure"
    "encoding/xml"
    "fmt"
    "net/http"
    "time"
)

func apiLatestHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, xml.Header)
    enc := xml.NewEncoder(w)
    feeds, _ := db.GetAllFeeds()
    entries := structure.ExtractEntriesFromFeeds(feeds...)
    filter, _ := db.GetFilter()
    entries = structure.FilterEntries(entries, filter)
    atomfeed := atom.CreateAtomFeed("Latest", "The latest news from your Brooklet subscriptions", "latestbrooklet", "Brooklet", "", time.Now())
    for _, entry := range entries {
        atomfeed.AddEntry(entry.Title, entry.Summary, entry.Content, entry.ID, entry.ID, entry.PermaLink, entry.Author, entry.Twitter, entry.Updated)
    }
    enc.Encode(atomfeed)
}

func apiFeedHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, xml.Header)
    enc := xml.NewEncoder(w)
    vars := mux.Vars(r)
    name := vars["name"]
    feed, _ := db.GetFeedByName(name)
    entries := structure.ExtractEntriesFromFeeds(feed)
    filter, _ := db.GetFilter()
    entries = structure.FilterEntries(entries, filter)
    atomfeed := atom.CreateAtomFeed(feed.Title, feed.Subtitle, feed.Name, feed.Title, feed.Twitter, feed.Updated)
    for _, entry := range entries {
        atomfeed.AddEntry(entry.Title, entry.Summary, entry.Content, entry.ID, entry.ID, entry.PermaLink, entry.Author, entry.Twitter, entry.Updated)
    }
    enc.Encode(atomfeed)
}