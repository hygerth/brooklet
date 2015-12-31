package brooklet

import (
    "github.com/gorilla/mux"
    "github.com/hygerth/brooklet/skywalker"
    "github.com/hygerth/brooklet/utils"
    "log"
    "net/http"
)

func Start() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    skywalker.Start()
    path, _ := utils.GetPath()
    r := mux.NewRouter()
    r.HandleFunc("/", homeHandler)
    r.HandleFunc("/latest", latestHandler)
    r.HandleFunc("/settings", settingsHandler)
    r.HandleFunc("/add/feed", addFeedHandler).Methods("POST")
    r.HandleFunc("/remove/feed", removeFeedHandler).Methods("POST")
    r.HandleFunc("/feed/{name}", feedHandler).Methods("GET")
    r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(path + "/static/"))))
    r.PathPrefix("/images").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir(path + "/images/"))))
    log.Println("Port: 9876")
    panic(http.ListenAndServe(":9876", r))
}