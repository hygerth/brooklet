package brooklet

import (
    "github.com/gorilla/mux"
    "github.com/hygerth/brooklet/db"
    "github.com/hygerth/brooklet/utils"
    "encoding/xml"
    "fmt"
    "net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello World")
}

func addFeedHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    utils.Checkerr(err)
    feed := r.Form["feed"][0]
    newfeed, _ := db.AddFeed(feed)
    http.Redirect(w, r, "/feed/" + newfeed.Name, 302)
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]
    feed, _ := db.GetFeedByName(name)
    xml.NewEncoder(w).Encode(feed)
}