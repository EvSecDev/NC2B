package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

/* moz_cookies all columns+types
id                        int
originAttributes          string
name                      string
value                     string
host                      string
path                      string
expiry                    int
lastAccessed              int
creationTime              int
isSecure                  int
isHttpOnly                int
inBrowserElement          int
sameSite                  int
schemeMap                 int
isPartitionedAttributeSet int
*/

func writeCookiesToFirefox(cookies []cookie, dbPath string) (err error) {
	timeout := 2000 // Timeout in milliseconds
	connStr := fmt.Sprintf("file:%s?_timeout=%d", dbPath, timeout)

	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		err = fmt.Errorf("error opening database at '%s': %v", dbPath, err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		err = fmt.Errorf("error connecting to database: %v", err)
		return
	}

	query := `
		INSERT OR REPLACE INTO moz_cookies (
			host,
			path,
			isSecure,
			isHttpOnly,
			expiry,
			name,
			value,
			originAttributes, 
			lastAccessed,
			creationTime,
			inBrowserElement,
			sameSite,
			schemeMap,
			isPartitionedAttributeSet
		)
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			'', -- originAttributes, assuming empty string as placeholder
			?,  -- lastAccessed
			?,  -- creationTime
			0,  -- inBrowserElement, assuming 0
			0,  -- sameSite, assuming 0
			0,  -- schemeMap, assuming 0
			0   -- isPartitionedAttributeSet, assuming 0
		);
	`

	for _, cookie := range cookies {
		// Convert boolean to integer, 1=true
		isSecureInt := 0
		if cookie.isSecure {
			isSecureInt = 1
		}
		isHttpOnlyInt := 0
		if cookie.httpOnly {
			isHttpOnlyInt = 1
		}

		// Handle marking cookies valid for all subdomains
		if cookie.includeSub && !strings.HasPrefix(cookie.domain, ".") {
			cookie.domain = "." + cookie.domain
		}

		// Using import time as create/last access
		creationTime := time.Now().UnixMicro()
		lastAccessTime := creationTime

		_, err = db.Exec(query, cookie.domain, cookie.path, isSecureInt, isHttpOnlyInt, cookie.expiration, cookie.name, cookie.value, lastAccessTime, creationTime)
		if err != nil {
			err = fmt.Errorf("error updating cookie %s (%s): %v", cookie.name, cookie.domain, err)
			return
		}
	}

	return
}
