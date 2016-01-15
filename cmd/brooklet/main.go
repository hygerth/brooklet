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
    notes            = "Important notes: brooklet expects MongoDB to run locally on port 27017"
    options          = "Options:\n-h, -help\tPrint this help text and exit \n-v, -version\tPrint program version and exit\n" + cacheduration + "\n" + cachetimeunit + "\n" + serverport + "\n"
    version          = "2016.01.15"
    help             = fmt.Sprintf("%s\n%s\nVersion: %s\n%s", usage, notes, version, options)
    cacheduration    = "-duration\tSpecify duration for cached content, default 5"
    cachetimeunit    = "-timeunit\tSpecify time unit for cache duration, default minutes"
    serverport       = "-port\t\tSpecify server port, default 9876"
    cliVersion       = flag.Bool("version", false, version)
    cliHelp          = flag.Bool("help", false, help)
    cliCacheDuration = flag.String("duration", "5", cacheduration)
    cliCacheTimeUnit = flag.String("timeunit", "minutes", cachetimeunit)
    cliServerPort    = flag.String("port", "9876", serverport)
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
    optionssettings["serverport"] = *cliServerPort

    brooklet.Start(optionssettings)
    exit(0)
    return
}
