package main

type browserInfo struct {
	name        string
	dbTable     string
	handlerFunc func([]cookie, string, string) (err error)
}

func setupBrowserSettings() (browsers map[string]browserInfo) {
	browsers = make(map[string]browserInfo)

	ffBrowser := browserInfo{
		name:        "firefox",
		dbTable:     "moz_cookies",
		handlerFunc: writeCookiesToFirefox,
	}
	browsers[ffBrowser.name] = ffBrowser

	chromeBrowser := browserInfo{
		name:        "chrome",
		dbTable:     "cookies",
		handlerFunc: writeCookiesToChrome,
	}
	browsers[chromeBrowser.name] = chromeBrowser

	return
}
