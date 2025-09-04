package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type cookie struct {
	domain     string
	includeSub bool
	path       string
	isSecure   bool
	httpOnly   bool
	expiration int
	name       string
	value      string
}

func parseNetscapeCookies(cookieData []byte) (parsedCookies []cookie, err error) {
	const httpOnlyPrefix string = "#HttpOnly_"
	// Create a scanner to read the cookie file from bytes
	scanner := bufio.NewScanner(bytes.NewReader(cookieData))

	for scanner.Scan() {
		line := scanner.Text()

		// Skipping lines to ignore
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, httpOnlyPrefix) {
			continue
		}

		fields := strings.Fields(line)

		if len(fields) != 7 {
			continue
		}

		var expiryEpoch int
		expiryEpoch, err = strconv.Atoi(fields[4])
		if err != nil {
			err = fmt.Errorf("found invalid epoch number in cookie %s (%s): %v", fields[5], fields[0], err)
			return
		}

		var cookieHTTPOnly bool
		var cookieDomain string
		if strings.HasPrefix(fields[0], httpOnlyPrefix) {
			cookieHTTPOnly = true
			cookieDomain = strings.TrimPrefix(string(fields[0]), httpOnlyPrefix)
		} else {
			cookieDomain = fields[0]
		}

		var cookieSecureOnly bool
		switch fields[3] {
		case "TRUE":
			cookieSecureOnly = true
		case "FALSE":
			cookieSecureOnly = false
		default:
			err = fmt.Errorf("found invalid isSecure field in cookie %s (%s): %v", fields[5], fields[0], err)
			return
		}

		var cookieIncSubDomains bool
		switch fields[1] {
		case "TRUE":
			cookieIncSubDomains = true
		case "FALSE":
			cookieIncSubDomains = false
		default:
			err = fmt.Errorf("found invalid includeSubdomains field in cookie %s (%s): %v", fields[5], fields[0], err)
			return
		}

		cookie := cookie{
			domain:     cookieDomain,
			includeSub: cookieIncSubDomains,
			path:       fields[2],
			isSecure:   cookieSecureOnly,
			httpOnly:   cookieHTTPOnly,
			expiration: expiryEpoch,
			name:       fields[5],
			value:      fields[6],
		}

		// Add the cookie to the list
		parsedCookies = append(parsedCookies, cookie)
	}
	err = scanner.Err()
	if err != nil {
		return
	}

	if len(parsedCookies) == 0 {
		err = fmt.Errorf("no valid cookies found")
		return
	}

	return
}
