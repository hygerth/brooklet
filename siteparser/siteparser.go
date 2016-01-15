package siteparser

import (
    "bytes"
    "github.com/hygerth/brooklet/utils"
    "golang.org/x/net/html"
    "regexp"
    "strings"
)

var iconregex = regexp.MustCompile(`icon[^>]*href=['|"]([^'"]+)['|"]`)
var twittersiteregex = regexp.MustCompile(`twitter:site['|"] content=['|"]([^'"]+)['|"]`)

var twittercreatorregex = regexp.MustCompile(`twitter:creator['|"] content=['|"]([^'"]+)['|"]`)
var ogdescriptionregex = regexp.MustCompile(`og:description['|"] content=['|"]([^'"]+)['|"]`)
var ogimageregex = regexp.MustCompile(`og:image['|"] content=['|"]([^'"]+)['|"]`)
var ogimageregex2 = regexp.MustCompile(`content=['|"]([^'"]+)['|"][^>]*og:image`)

type Meta struct {
    Icon           string
    TwitterSite    string
    TwitterCreator string
    Description    string
    Image          string
}

var positivecandidatesregex = regexp.MustCompile(`and|article|body|column|main`)
var negativecandidatesregex = regexp.MustCompile(`combx|modal|lightbox|comment|disqus|foot|footer|head|header|menu|meta|nav|rss|script|shoutbox|sidebar|sponsor|social|teaserlist|time|tweet|twitter`)

var positiveregex = regexp.MustCompile(`article|body|content|entry|hentry|page|pagination|post|text`)
var negativeregex = regexp.MustCompile(`ad|banner|brand|combx|comment|comments|contact|foot|footer|footnote|left|link|media|meta|navigation|promo|related|right|scroll|share|shoutbox|sidebar|sponsor|utility|tags|widget`)

var containerregrex = regexp.MustCompile(`article|aside|div|section`)

var nonwordregex = regexp.MustCompile(`[^A-Za-z0-9]+`)

func GetMetaForSite(url string) (Meta, error) {
    page, err := utils.GetPage(url)
    if err != nil {
        return Meta{}, err
    }
    meta, err := extractMetaFromSiteData(page, url)
    if err != nil {
        return Meta{}, err
    }
    return meta, nil
}

func extractMetaFromSiteData(data []byte, url string) (Meta, error) {
    var meta Meta
    str := string(data)
    found := iconregex.FindAllStringSubmatch(str, -1)
    for _, icon := range found {
        if strings.Contains(icon[1], "favicon") {
            img, err := utils.RelativeToAbsolutePath(icon[1], url)
            if err != nil {
                return meta, err
            }
            meta.Icon = img
        }
    }
    f := twittersiteregex.FindStringSubmatch(str)
    if len(f) > 1 {
        meta.TwitterSite = f[1]
    }

    f = twittercreatorregex.FindStringSubmatch(str)
    if len(f) > 1 {
        meta.TwitterCreator = f[1]
    }
    f = ogdescriptionregex.FindStringSubmatch(str)
    if len(f) > 1 {
        meta.Description = f[1]
    }
    f = ogimageregex.FindStringSubmatch(str)
    if len(f) > 1 {
        img, err := utils.RelativeToAbsolutePath(f[1], url)
        if err != nil {
            return meta, err
        }
        meta.Image = img
    } else {
        f = ogimageregex2.FindStringSubmatch(str)
        if len(f) > 1 {
            img, err := utils.RelativeToAbsolutePath(f[1], url)
            if err != nil {
                return meta, err
            }
            meta.Image = img
        }
    }
    return meta, nil
}

// GetArticleForSite is not guaranteed to give good results for all URLs
func GetArticleForSite(url string) (string, error) {
    page, err := utils.GetPage(url)
    if err != nil {
        return "", err
    }
    article := getArticle(page)
    return article, nil
}

func getArticle(data []byte) string {
    r := bytes.NewReader(data)
    doc, _ := html.Parse(r)
    // Tags
    doc = removeNegativeCandidates(doc)
    doc = removeNegativeMatches(doc)
    doc = getBodyElement(doc)
    // Attributes
    doc = removeNegativeAttributeMatches(doc)
    doc, _ = retriveMainRole(doc)
    doc = removeNonMainContent(doc)
    doc = clearClassesAndIDs(doc)
    c := calcContent(doc)
    var buff bytes.Buffer
    html.Render(&buff, doc)
    articlestr := buff.String()
    articlestr = utils.RemoveNewLines(articlestr)
    articlestr = utils.ReplaceTabsWithASpace(articlestr)
    articlestr = utils.TrimSpaces(articlestr)
    if float64(c)/float64(len(articlestr)) < 0.2 {
        // At least 20% of the article should be text
        return ""
    }
    return articlestr
}

func calcContent(n *html.Node) int {
    counter := 0
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        if c.Type == html.TextNode {
            a := nonwordregex.ReplaceAllString(c.Data, "")
            counter += len(a)
        } else {
            counter += calcContent(c)
        }
    }
    return counter
}

func removeNegativeCandidates(n *html.Node) *html.Node {
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        if !negativecandidatesregex.MatchString(c.Data) || positivecandidatesregex.MatchString(c.Data) {
            d := removeNegativeCandidates(c)
            if c.PrevSibling != nil {
                c.PrevSibling.NextSibling = d
            } else {
                n.FirstChild = d
            }
        } else {
            if c.PrevSibling != nil {
                c.PrevSibling.NextSibling = c.NextSibling
            } else {
                n.FirstChild = c.NextSibling
            }
        }
    }
    return n
}

func removeNegativeMatches(n *html.Node) *html.Node {
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        if !negativeregex.MatchString(c.Data) || positiveregex.MatchString(c.Data) {
            d := removeNegativeMatches(c)
            if c.PrevSibling != nil {
                c.PrevSibling.NextSibling = d
            } else {
                n.FirstChild = d
            }
        } else {
            if c.PrevSibling != nil {
                c.PrevSibling.NextSibling = c.NextSibling
            } else {
                n.FirstChild = c.NextSibling
            }
        }
    }
    return n
}

func getBodyElement(n *html.Node) *html.Node {
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        if strings.ToLower(c.Data) == "body" {
            return c
        } else {
            d := getBodyElement(c)
            if strings.ToLower(d.Data) == "body" {
                return d
            }
        }
    }
    return n
}

func retriveMainRole(n *html.Node) (*html.Node, bool) {
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        for _, attr := range c.Attr {
            key := strings.ToLower(attr.Key)
            val := strings.ToLower(attr.Val)
            if key == "role" && val == "main" {
                return c, true
            }
        }
        d, hasMainRole := retriveMainRole(c)
        if hasMainRole {
            return d, hasMainRole
        }
    }
    return n, false
}

func removeNonMainContent(n *html.Node) *html.Node {
    var max int
    maincontent := n
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        value := calculateContentValue(c)
        if value > max {
            maincontent = c
            max = value
        }
    }
    return maincontent
}

func calculateContentValue(n *html.Node) int {
    var value int
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        if c.Type != html.TextNode {
            value += calculateContentValue(c)
        } else {
            value += len(c.Data)
        }
    }
    return value
}

func removeNegativeAttributeMatches(n *html.Node) *html.Node {
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        if c.Type != html.TextNode && containerregrex.MatchString(c.Data) {
            for _, attr := range c.Attr {
                key := strings.ToLower(attr.Key)
                if key == "id" || key == "class" {
                    val := strings.ToLower(attr.Val)
                    values := nonwordregex.Split(val, -1)
                    penalty := 0
                    for _, value := range values {
                        if negativeregex.MatchString(value) {
                            penalty = penalty + 4
                        }
                    }
                    if penalty > 0 {
                        if c.PrevSibling != nil {
                            c.PrevSibling.NextSibling = c.NextSibling
                        } else {
                            n.FirstChild = c.NextSibling
                        }
                    } else {
                        d := removeNegativeAttributeMatches(c)
                        if c.PrevSibling != nil {
                            c.PrevSibling.NextSibling = d
                        } else {
                            n.FirstChild = c.NextSibling
                        }
                    }
                }
            }
        }
    }
    return n
}

func clearClassesAndIDs(n *html.Node) *html.Node {
    for i, attr := range n.Attr {
        key := strings.ToLower(attr.Key)
        if key == "id" || key == "class" {
            n.Attr[i].Val = ""
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        c = clearClassesAndIDs(c)
    }
    return n
}
