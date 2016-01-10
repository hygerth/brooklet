package skywalker

import (
    "github.com/hygerth/brooklet/db"
    "github.com/hygerth/brooklet/structure"
    "github.com/hygerth/brooklet/utils"
    "log"
    "time"
)

// Start starts the background update schedule of Skywalker
func Start() {
    Sync()
    go func() {
        timer := time.Tick(10 * time.Minute)
        for _ = range timer {
            log.Println("skywalker: Syncing feeds")
            Sync()
            log.Println("skywalker: Syncing completed")
        }
    }()
}

func Sync() {
    var feedsToUpdate []structure.Feed
    feeds, err := db.GetAllFeeds()
    utils.Checkerr(err)
    for _, feed := range feeds {
        maxTimeSinceUpdate := 24 / feed.ChangeFrequency
        if time.Now().Sub(feed.LastUpdate).Hours() > maxTimeSinceUpdate {
            feedsToUpdate = append(feedsToUpdate, feed)
        }
    }
    processFeeds(feedsToUpdate)
}