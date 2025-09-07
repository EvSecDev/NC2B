package main

import (
	"fmt"
	"strings"
	"time"
)

func filterCookies(allCookies []cookie, filters cookieFilters) (filteredCookies []cookie, err error) {
	for _, cookie := range allCookies {
		// Ignore expired cookies, when expiration is present
		currentTime := time.Now().Unix()
		if currentTime > int64(cookie.expiration) && cookie.expiration > 0 && !filters.allowExpired {
			continue
		}

		// For any provided filters, skip non-matching cookie fields
		if filters.name != "" && !strings.Contains(cookie.name, filters.name) {
			continue
		}
		if filters.domain != "" && !strings.HasSuffix(cookie.domain, filters.domain) {
			continue
		}

		// Add to filtered list
		filteredCookies = append(filteredCookies, cookie)
	}

	if len(filteredCookies) == 0 {
		err = fmt.Errorf("no cookies found for search parameters")
		return
	}
	return
}
