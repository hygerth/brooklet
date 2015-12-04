package skywalker

import (
    "github.com/hygerth/brooklet/utils"
    "fmt"
    "image"
    _ "image/gif"
    _ "image/jpeg"
    _ "image/png"
    "io"
    "net/http"
    "os"
    "os/exec"
    "regexp"
)

const useragent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.3"

var filenameregex = regexp.MustCompile(`[^/]/([^/]+\.[a-z]{2,4})`)

// SyncImage, return base-filename, IsPortrait and error
func SyncImage(url string) (string, bool, error) {
    if len(url) == 0 {
        return "", false, nil
    }
    sizes := []string{"128", "512"}
    path, err := utils.GetPath()
    if err != nil {
        return "", false, err
    }
    path += "/images/"
    os.Mkdir(path, os.ModePerm)
    filename := getFilename(url)
    picture, err := downloadPicture(path, filename, url)
    if err != nil {
        return "", false, err
    }
    img := picture.Bounds()
    height := img.Dy()
    width := img.Dx()
    portrait := false
    if height > width {
        portrait = true
    }
    uniquefilename, err := utils.GenerateUniqueFilename(path)
    if err != nil {
        return filename, portrait, err
    }
    for _, size := range sizes {
        err = resizeImage(path, filename, uniquefilename, size)
        if err != nil {
            return uniquefilename, portrait, err
        }
    }
    err = os.Remove(path + filename)
    if err != nil {
        return uniquefilename, portrait, err
    }
    return uniquefilename, portrait, nil
}

func getFilename(url string) string {
    found := filenameregex.FindStringSubmatch(url)
    if len(found) > 1 {
        return found[1]
    }
    return "unknown.png"
}

func resizeImage(path, original, newimage, size string) error {
    newfile := path + newimage + "-" + size + ".png"
    args := []string{"-s", size, "-o", newfile, path + original}
    cmd := exec.Command("vipsthumbnail", args...)
    err := cmd.Start()
    if err != nil {
        return err
    }
    err = cmd.Wait()
    return err
}

func downloadPicture(path, filename, url string) (image.Image, error) {
    out, err := os.Create(path + filename)
    defer out.Close()
    if err != nil {
        return nil, err
    }
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("User-Agent", useragent)
    resp, err := client.Do(req)
    if err != nil || resp.StatusCode != http.StatusOK {
        if err == nil {
            err = fmt.Errorf("utils: Server responded with %d, %s", resp.StatusCode, resp.Status)
        }
        return nil, err
    }
    defer resp.Body.Close()
    io.Copy(out, resp.Body)
    out.Seek(0, 0)
    m, _, err := image.Decode(out)
    return m, err
}