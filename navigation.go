package brooklet

import(
    "github.com/hygerth/brooklet/db"
    "sort"
    "strings"
)

var navigationbase = map[string]string{
    "Home": "/",
    "Latest": "/latest",
    "Settings": "/settings",
}

type Navigation struct {
    NavigationItems []NavigationItem `xml:"navigationitem"`
}

type NavigationItem struct {
    Title string `xml:"title,attr"`
    URL string `xml:"url,attr"`
    SubList SubList`xml:"sublist"`
}

type NavigationItems []NavigationItem

type SubList struct {
    SubItems []SubItem `xml:"item"`
}

type SubItem struct {
    Title string `xml:"title"`
    URL string `xml:"url"`
    Icon string `xml:"icon"`
}

func buildNavigation() Navigation {
    var nav Navigation
    for key, _ := range navigationbase {
        navitem := NavigationItem{Title: key, URL: navigationbase[key]}
        nav.NavigationItems = append(nav.NavigationItems, navitem)
    }
    feeds, _ := db.GetAllFeeds()
    var sublist SubList
    for _, feed := range feeds {
        subitem := SubItem{Title: feed.Title, URL: "/feed/" + feed.Name, Icon: feed.Icon}
        sublist.SubItems = append(sublist.SubItems, subitem)
    }
    nav.NavigationItems = append(nav.NavigationItems, NavigationItem{Title: "Subscriptions", SubList: sublist})
    nav.NavigationItems = SortNavigationItems(nav.NavigationItems)
    return nav
}

func (nv NavigationItems) Len() int {
    return len(nv)
}

func (nv NavigationItems) Less(i, j int) bool {
    return strings.ToLower(nv[i].Title) < strings.ToLower(nv[j].Title)
}

func (nv NavigationItems) Swap(i, j int) {
    nv[i], nv[j] = nv[j], nv[i]
}

// SortEntriesByDate sorts the entries in a list by the date of which they
// were updated on
func SortNavigationItems(nv []NavigationItem) []NavigationItem {
    nvsorted := make(NavigationItems, 0, len(nv))
    for _, nvi := range nv {
        nvsorted = append(nvsorted, nvi)
    }
    sort.Sort(nvsorted)
    var items []NavigationItem
    for _, item := range nvsorted {
        items = append(items, item)
    }
    return items
}