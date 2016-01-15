package utils

import (
    "go/build"
    "fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "net/http"
    "net/url"
    "os"
    "regexp"
    "strconv"
    "strings"
    "time"
)

const useragent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.3"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var formats []string = []string{time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano, time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano}

func Checkerr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

// GetPath returns the path to the build directory
func GetPath() (string, error) {
    p, err := build.Default.Import("github.com/hygerth/brooklet", "", build.FindOnly)
    return p.Dir, err
}

func GetPage(url string) ([]byte, error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return []byte{}, err
    }
    req.Header.Set("User-Agent", useragent)
    resp, err := client.Do(req)
    if err != nil || resp.StatusCode != http.StatusOK {
        if err == nil {
            err = fmt.Errorf("utils: Server responded with %d, %s", resp.StatusCode, resp.Status)
        }
        return []byte{}, err
    }
    defer resp.Body.Close()
    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return []byte{}, err
    }
    return b, nil
}

func StringToInt64(s string) int64 {
    i, err := strconv.ParseInt(s, 10, 64)
    Checkerr(err)
    return i
}

func StringToTimeDurationUnit(s string) time.Duration {
    switch s {
    case "days": return 24 * time.Hour
    case "hours": return time.Hour
    default: return time.Minute
    }
}

func ParseTimeString(str string) (time.Time, error) {
    var t time.Time
    for _, format := range formats {
        t, err := time.Parse(format, str)
        if err == nil {
            return t, nil
        }
    }
    return t, fmt.Errorf("utils: Could not parse '%s' as time", str)
}

func RelativeToAbsolutePath(relative string, hosturl string) (string, error) {
    r, err := url.Parse(relative)
    if err != nil {
        return relative, err
    }
    if r.IsAbs() {
        return relative, nil
    }
    u, err := url.Parse(hosturl)
    if err != nil {
        return relative, err
    }
    r = r.ResolveReference(u)
    return r.String(), nil
}

func RemoveNewLines(s string) string {
    re := regexp.MustCompile(`\n`)
    return re.ReplaceAllString(s, " ")
}

func ReplaceTabsWithASpace(s string) string {
    re := regexp.MustCompile(`\t`)
    return re.ReplaceAllString(s, " ")
}

func TrimSpaces(s string) string {
    re := regexp.MustCompile(`\s{2,}`)
    return re.ReplaceAllString(s, " ")
}

func GenerateUniqueID() string {
    var src = rand.NewSource(time.Now().UnixNano())
    id := randStringBytesMaskImprSrc(16, src)
    return id
}

func GenerateUniqueFilename(path string) (string, error) {
    var src = rand.NewSource(time.Now().UnixNano())
    filename := randStringBytesMaskImprSrc(16, src)
    directory, err := os.Open(path)
    if err != nil {
        return filename, err
    }
    defer directory.Close()

    files, err := directory.Readdir(-1)
    if err != nil {
        return filename, err
    }

    for _, file := range files {
        if strings.Contains(file.Name(), filename) {
            directory.Close()
            return GenerateUniqueFilename(path)
        }
    }
    return filename, nil
}

// RandStringBytesMaskImprSrc source: http://stackoverflow.com/a/31832326
func randStringBytesMaskImprSrc(n int, src rand.Source) string {
    b := make([]byte, n)
    for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = src.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }

    return string(b)
}

func IsMobile(header http.Header) bool {
    var ua string
    uagets := []string{"user-agent", "User-Agent", "USER-AGENT"}
    for _, uaget := range uagets {
        if header.Get(uaget) != "" {
            ua = header.Get(uaget)
            break
        }
    }
    if ua == "" {
        return false
    }
    if strings.Contains(strings.ToLower(ua), "mobile") && !strings.Contains(strings.ToLower(ua), "ipad") {
        return true
    }
    return false
}