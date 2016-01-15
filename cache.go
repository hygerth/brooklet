package brooklet

import (
    "github.com/hygerth/brooklet/db"
    "github.com/hygerth/brooklet/structure"
    "github.com/hygerth/brooklet/utils"
    "log"
    "sync"
    "time"
)

type Cache struct {
    lock sync.RWMutex
    duration int64
    timeunit time.Duration
    data map[string]Data
}

type Data struct {
    timestamp time.Time
    Title string
    Entries []structure.Entry
}

func (c *Cache) Init(options map[string]string) {
    c.parseOptions(options)
    c.data = make(map[string]Data)
    go func() {
        timer := time.Tick(1 * time.Minute)
        for _ = range timer {
            //log.Println("brooklet: Cleaning cache")
            c.lock.RLock()
            max := time.Duration(c.duration) * c.timeunit
            for key := range c.data {
                if time.Since(c.data[key].timestamp) > max {
                    log.Printf("brooklet: Purged '%s' from cache\n", key)
                    delete(c.data, key)
                }
            }
            c.lock.RUnlock()
            //log.Println("brooklet: Cleaning cache completed")
        }
    }()
}

func (c *Cache) parseOptions(options map[string]string) {
    c.duration = utils.StringToInt64(options["cacheduration"])
    c.timeunit = utils.StringToTimeDurationUnit(options["cachetimeunit"])
}

func (c *Cache) Empty() {
    log.Println("brooklet: Emptying cache")
    c.lock.RLock()
    defer c.lock.RUnlock()
    for key := range c.data {
        delete(c.data, key)
    }
}

func (c *Cache) Get(key string) (*Data, bool) {
    c.lock.RLock()
    defer c.lock.RUnlock()
    d, ok := c.data[key]
    if !ok || time.Since(d.timestamp) > time.Duration(c.duration) * c.timeunit {
        c.lock.RUnlock()
        c.Set(key)
        c.lock.RLock()
        return c.Get(key)
    }
    log.Println("brooklet: Serving from cache")
    return &d, ok
}

func (c *Cache) Set(key string) {
    log.Printf("brooklet: Adding '%s' to cache\n", key)
    if key == "all" {
        c.SetAll()
        return
    }
    c.lock.Lock()
    defer c.lock.Unlock()
    feed, _ := db.GetFeedByName(key)
    entries := structure.ExtractEntriesFromFeeds(feed)
    filter, _ := db.GetFilter()
    entries = structure.FilterEntries(entries, filter)
    if len(entries) > 100 {
        entries = entries[:100]
    }
    var d Data
    d.timestamp = time.Now()
    d.Title = feed.Title
    d.Entries = entries
    c.data[key] = d
}

func (c *Cache) SetAll() {
    c.lock.Lock()
    defer c.lock.Unlock()
    feeds, _ := db.GetAllFeeds()
    entries := structure.ExtractEntriesFromFeeds(feeds...)
    filter, _ := db.GetFilter()
    entries = structure.FilterEntries(entries, filter)
    if len(entries) > 100 {
        entries = entries[:100]
    }
    var d Data
    d.timestamp = time.Now()
    d.Title = "All"
    d.Entries = entries
    c.data["all"] = d
}