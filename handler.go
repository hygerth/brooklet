package brooklet

import (
    "github.com/gorilla/mux"
    "github.com/hygerth/brooklet/db"
    "github.com/hygerth/brooklet/skywalker"
    "github.com/hygerth/brooklet/structure"
    "github.com/hygerth/brooklet/utils"
    "encoding/xml"
    "fmt"
    "net/http"
    "strings"
)

const xslLocation = "/static/xsl/"

var viewheaders = map[string]string{
    "feed": "<?xml version=\"1.0\"?><?xml-stylesheet type=\"text/xsl\" href=\"_xslfile_\"?><!DOCTYPE page [<!ENTITY % page SYSTEM \"/static/dtd/page.dtd\">%page;<!ENTITY % content SYSTEM \"/static/dtd/content.dtd\">%content;<!ENTITY % navigation SYSTEM \"/static/dtd/navigation.dtd\">%navigation;]>",
    "home": "<?xml version=\"1.0\"?><?xml-stylesheet type=\"text/xsl\" href=\"_xslfile_\"?>",
    "settings": "<?xml version=\"1.0\"?><?xml-stylesheet type=\"text/xsl\" href=\"_xslfile_\"?>",
}

type Page struct {
    XMLName xml.Name `xml:"page"`
    Navigation Navigation `xml:"navigation"`
    Content Content `xml:"content"`
    Subscriptions []Subscription `xml:"subscription"`
}

type Content struct {
    Title string `xml:"title,attr"`
    Entries []structure.Entry `xml:"entry"`
}

type Subscription struct {
    Title string `xml:"title"`
    URL string `xml:"url"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    ismobile := utils.IsMobile(r.Header)
    xmlheader, err := buildXMLHeader("home", ismobile)
    utils.Checkerr(err)
    fmt.Fprintf(w, "%s", xmlheader)
    enc := xml.NewEncoder(w)
    p := buildBasicPage()
    enc.Encode(p)
}

func latestHandler(w http.ResponseWriter, r *http.Request) {
    ismobile := utils.IsMobile(r.Header)
    xmlheader, err := buildXMLHeader("feed", ismobile)
    utils.Checkerr(err)
    fmt.Fprintf(w, "%s", xmlheader)
    enc := xml.NewEncoder(w)
    var p Page
    nav := buildNavigation()
    p.Navigation = nav
    feeds, _ := db.GetAllFeeds()
    entries := structure.ExtractEntriesFromFeeds(feeds...)
    content := Content{Title: "Latest", Entries: entries}
    p.Content = content
    enc.Encode(p)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
    ismobile := utils.IsMobile(r.Header)
    xmlheader, err := buildXMLHeader("settings", ismobile)
    utils.Checkerr(err)
    fmt.Fprintf(w, "%s", xmlheader)
    enc := xml.NewEncoder(w)
    p := buildBasicPage()
    enc.Encode(p)
}

func buildBasicPage() Page {
    var p Page
    nav := buildNavigation()
    p.Navigation = nav
    subscriptions := buildSubscriptions()
    p.Subscriptions = subscriptions
    return p
}

func addFeedHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    utils.Checkerr(err)
    feed := r.Form["feed"][0]
    newfeed, _ := db.AddFeed(feed)
    skywalker.UpdateFeed(newfeed)
    http.Redirect(w, r, "/feed/" + newfeed.Name, http.StatusFound)
}

func removeFeedHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    utils.Checkerr(err)
    feed := r.Form["feed"][0]
    err = db.RemoveFeed(feed)
    utils.Checkerr(err)
    http.Redirect(w, r, "/", http.StatusFound)
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
    ismobile := utils.IsMobile(r.Header)
    xmlheader, err := buildXMLHeader("feed", ismobile)
    utils.Checkerr(err)
    fmt.Fprintf(w, "%s", xmlheader)
    enc := xml.NewEncoder(w)
    vars := mux.Vars(r)
    name := vars["name"]
    var p Page
    nav := buildNavigation()
    p.Navigation = nav
    feed, _ := db.GetFeedByName(name)
    content := Content{Title: feed.Title, Entries: structure.ExtractEntriesFromFeeds(feed)}
    p.Content = content
    enc.Encode(p)
}

func buildSubscriptions() []Subscription {
    var subscriptions []Subscription
    feeds, _ := db.GetAllFeeds()
    for _, feed := range feeds {
        subscription := Subscription{Title: feed.Title, URL: feed.URL}
        subscriptions = append(subscriptions, subscription)
    }
    return subscriptions
}

func buildXMLHeader(view string, isMobile bool) (string, error) {
    i, ok := viewheaders[view]
    if !ok {
        return "", fmt.Errorf("brooklet: Could not find a header for view '%s'", view)
    }
    if isMobile {
        view = view + "-mobile"
    }
    return strings.Replace(i, "_xslfile_", xslLocation + view + ".xsl", 1), nil
}
