package main

import (
    "flag"
    "fmt"
    "github.com/hygerth/brooklet"
    "os"
)

var exit = os.Exit
var optionssettings map[string]string = make(map[string]string)

var (
    usage            = "Usage: brooklet [OPTIONS]"
    options          = "Options:\n-h, -help\tPrint this help text and exit \n-v, -version\tPrint program version and exit\n" + cacheduration + "\n" + cachetimeunit + "\n"
    version          = "2016.01.11"
    help             = fmt.Sprintf("%s\nVersion: %s\n%s", usage, version, options)
    cacheduration    = "-duration\tSpecify duration for cached content, default 5"
    cachetimeunit    = "-timeunit\tSpecify time unit for cache duration, default minutes"
    cliVersion       = flag.Bool("version", false, version)
    cliHelp          = flag.Bool("help", false, help)
    cliCacheDuration = flag.String("duration", "5", cacheduration)
    cliCacheTimeUnit = flag.String("timeunit", "minutes", cachetimeunit)
)

func init() {
    flag.BoolVar(cliVersion, "v", false, version)
    flag.BoolVar(cliHelp, "h", false, help)
}

func main() {
    flag.Parse()

    if *cliVersion {
        fmt.Println(flag.Lookup("version").Usage)
        exit(0)
        return
    }
    if *cliHelp {
        fmt.Println(flag.Lookup("help").Usage)
        exit(0)
        return
    }
    optionssettings["cacheduration"] = *cliCacheDuration
    optionssettings["cachetimeunit"] = *cliCacheTimeUnit

    brooklet.Start(optionssettings)
    exit(0)
    return
}
