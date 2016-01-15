package brooklet

import (
    "github.com/gorilla/mux"
    "github.com/hygerth/brooklet/skywalker"
    "github.com/hygerth/brooklet/utils"
    "log"
    "net/http"
)

var cache Cache

func Start(options map[string]string) {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    skywalker.Start()
    cache = Cache{}
    cache.Init(options)
    path, _ := utils.GetPath()
    r := mux.NewRouter()
    r.HandleFunc("/", homeHandler).Methods("GET")
    r.HandleFunc("/latest", latestHandler).Methods("GET")
    r.HandleFunc("/settings", settingsHandler).Methods("GET")
    r.HandleFunc("/add/feed", addFeedHandler).Methods("POST")
    r.HandleFunc("/remove/feed", removeFeedHandler).Methods("POST")
    r.HandleFunc("/add/filter", addFilterHandler).Methods("POST")
    r.HandleFunc("/remove/filter", removeFilterHandler).Methods("POST")
    r.HandleFunc("/article/{id}", articleHandler).Methods("GET")
    r.HandleFunc("/feed/{name}", feedHandler).Methods("GET")
    r.HandleFunc("/api/latest", apiLatestHandler).Methods("GET")
    r.HandleFunc("/api/feed/{name}", apiFeedHandler).Methods("GET")
    r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(path + "/static/"))))
    r.PathPrefix("/images").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir(path + "/images/"))))
    log.Println("Port: " + options["serverport"])
    panic(http.ListenAndServe(":" + options["serverport"], r))
}