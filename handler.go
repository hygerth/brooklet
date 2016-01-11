package brooklet

import (
    "github.com/gorilla/mux"
    "github.com/hygerth/brooklet/db"
    "github.com/hygerth/brooklet/skywalker"
    "github.com/hygerth/brooklet/structure"
    "github.com/hygerth/brooklet/utils"
    "encoding/xml"
    "fmt"
    "html/template"
    "net/http"
    "strings"
    "time"
)

const xslLocation = "/static/xsl/"

var viewheaders = map[string]string{
    "feed": "<?xml-stylesheet type=\"text/xsl\" href=\"_xslfile_\"?><!DOCTYPE page [<!ENTITY % page SYSTEM \"/static/dtd/page.dtd\">%page;<!ENTITY % content SYSTEM \"/static/dtd/content.dtd\">%content;<!ENTITY % navigation SYSTEM \"/static/dtd/navigation.dtd\">%navigation;]>",
    "home": "<?xml-stylesheet type=\"text/xsl\" href=\"_xslfile_\"?>",
    "settings": "<?xml-stylesheet type=\"text/xsl\" href=\"_xslfile_\"?>",
}

type Page struct {
    XMLName xml.Name `xml:"page"`
    Navigation Navigation `xml:"navigation,omitempty"`
    Content Content `xml:"content,omitempty"`
    Subscriptions []Subscription `xml:"subscription,omitempty"`
    Filter []string `xml:"filter,omitempty"`
}

type Content struct {
    Title string `xml:"title,attr"`
    Entries []structure.Entry `xml:"entry,omitempty"`
    Article structure.Entry `xml:"article,omitempty"`
}

type Subscription struct {
    Title string `xml:"title,omitempty"`
    ID string `xml:"id,omitempty"`
    URL string `xml:"url,omitempty"`
}

type Article struct {
    Article interface{}
    Content interface{}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, xml.Header)
    ismobile := utils.IsMobile(r.Header)
    xmlheader, err := buildXMLHeader("home", ismobile)
    utils.Checkerr(err)
    fmt.Fprintf(w, "%s", xmlheader)
    enc := xml.NewEncoder(w)
    p := buildBasicPage()
    enc.Encode(p)
}

func latestHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, xml.Header)
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
    filter, _ := db.GetFilter()
    entries = structure.FilterEntries(entries, filter)
    content := Content{Title: "Latest", Entries: entries[:20]}
    p.Content = content
    enc.Encode(p)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, xml.Header)
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
    filter := buildFilter()
    p.Filter = filter
    return p
}

func addFeedHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    utils.Checkerr(err)
    feed := r.Form["feed"][0]
    go func(feed string) {
        newfeed, _ := db.AddFeed(feed)
        skywalker.UpdateFeed(newfeed)
    }(feed)
    time.Sleep(4 * time.Second)
    http.Redirect(w, r, "/settings", http.StatusFound)
}

func removeFeedHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    utils.Checkerr(err)
    feed := r.Form["feed"][0]
    err = db.RemoveFeed(feed)
    utils.Checkerr(err)
    http.Redirect(w, r, "/settings", http.StatusFound)
}

func addFilterHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    utils.Checkerr(err)
    filter := r.Form["filter"][0]
    db.AddFilter(filter)
    http.Redirect(w, r, "/settings", http.StatusFound)
}

func removeFilterHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    utils.Checkerr(err)
    filter := r.Form["filter"][0]
    err = db.RemoveFilter(filter)
    utils.Checkerr(err)
    http.Redirect(w, r, "/settings", http.StatusFound)
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    article, _ := db.GetEntryByArticleID(id)
    c := template.HTML(article.Content)
    p := &Article{Article: article, Content: c}
    path, _ := utils.GetPath()
    t := template.Must(template.New("article.tmpl").ParseFiles(path + "/layouts/article.tmpl"))
    t.Execute(w, p)
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, xml.Header)
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
    entries := structure.ExtractEntriesFromFeeds(feed)
    filter, _ := db.GetFilter()
    entries = structure.FilterEntries(entries, filter)
    content := Content{Title: feed.Title, Entries: entries[:20]}
    p.Content = content
    enc.Encode(p)
}

func buildSubscriptions() []Subscription {
    var subscriptions []Subscription
    feeds, _ := db.GetAllFeeds()
    for _, feed := range feeds {
        subscription := Subscription{Title: feed.Title, ID: feed.Name, URL: feed.URL}
        subscriptions = append(subscriptions, subscription)
    }
    return subscriptions
}

func buildFilter() []string {
    var filter []string
    allfilters, _ := db.GetFilter()
    for _, f := range allfilters {
        filter = append(filter, f.Filter)
    }
    return filter
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
