package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
)

type cookieFilters struct {
	name         string
	domain       string
	allowExpired bool
}

func main() {
	var netscapeCookieFile string
	var browserName string
	var browserCookieStore string
	var filters cookieFilters
	var versionInfoRequested bool
	var versionRequested bool

	const usage = `
Netscape Cookie to Browser (NC2B)
  Retrieves cookies from netscape text files and injects into browser cookie storage

  Options:
    -j, --cookie-jar         <path/to/file>   Path to netscape cookie file to read cookies
    -b, --browser            <firefox|chrome> Browser name
    -c, --browser-cookie-jar <path/to/file>   Path to browser cookie storage file
    -n, --name               <text>           Import cookies matching name
    -d, --domain             <text>           Import cookies matching hostname
        --allow-expired                       Enable importing expired cookies
    -h, --help                                Show this help menu
    -V, --version                             Show version and packages
        --versionid                           Show only version number

Report bugs to: dev@evsec.net
NC2B home page: <https://github.com/EvSecDev/NC2B>
General help using GNU software: <https://www.gnu.org/gethelp/>
`
	flag.StringVar(&netscapeCookieFile, "j", "", "")
	flag.StringVar(&netscapeCookieFile, "cookie-jar", "", "")
	flag.StringVar(&browserName, "b", "", "")
	flag.StringVar(&browserName, "browser", "", "")
	flag.StringVar(&browserCookieStore, "c", "", "")
	flag.StringVar(&browserCookieStore, "browser-cookie-jar", "", "")
	flag.StringVar(&filters.name, "n", "", "")
	flag.StringVar(&filters.name, "name", "", "")
	flag.StringVar(&filters.domain, "d", "", "")
	flag.StringVar(&filters.domain, "domain", "", "")
	flag.BoolVar(&filters.allowExpired, "allow-expired", false, "")
	flag.BoolVar(&versionInfoRequested, "V", false, "")
	flag.BoolVar(&versionInfoRequested, "version", false, "")
	flag.BoolVar(&versionRequested, "versionid", false, "")

	flag.Usage = func() { fmt.Printf("Usage: %s [OPTIONS]...%s", os.Args[0], usage) }
	flag.Parse()

	const progVersion string = "v0.2.0"
	if versionInfoRequested {
		fmt.Printf("NC2B %s\n", progVersion)
		fmt.Printf("Built using %s(%s) for %s on %s\n", runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
		fmt.Print("Direct Package Imports: runtime strings strconv bufio _github.com/mattn/go-sqlite3 flag database/sql fmt time os bytes\n")
		return
	} else if versionRequested {
		fmt.Println(progVersion)
		return
	}

	if netscapeCookieFile == "" || browserCookieStore == "" || browserName == "" {
		fmt.Printf("No arguments specified or incorrect argument combination. Use '-h' or '--help' to guide your way.\n")
		return
	}

	// Config of browser settings
	browsers := setupBrowserSettings()

	// Validate user browser choice
	browserName = strings.TrimSpace(browserName)
	browserName = strings.ToLower(browserName)
	_, validBrowserChoice := browsers[browserName]
	if !validBrowserChoice {
		logError("Invalid option", fmt.Errorf("unsupported browser '%s'", browserName))
	}

	netscapeFileRaw, err := os.ReadFile(netscapeCookieFile)
	logError("Failed to read cookie jar", err)

	cookieList, err := parseNetscapeCookies(netscapeFileRaw)
	logError("Failed to parse cookie jar", err)

	cookiesToImport, err := filterCookies(cookieList, filters)
	logError("Failed to filter cookies", err)

	err = browsers[browserName].handlerFunc(cookiesToImport, browserCookieStore, browsers[browserName].dbTable)
	logError("Failed to import cookies into browser", err)
}

func logError(errorDescription string, errorMessage error) {
	if errorMessage == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "%s: %v\n", errorDescription, errorMessage)
	os.Exit(1)
}
