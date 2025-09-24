package main

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

/* table: moz_cookies all columns+types
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

func writeCookiesToFirefox(cookies []cookie, dbPath string, tableName string) (err error) {
	db, err := connectToDB(dbPath)
	if err != nil {
		err = fmt.Errorf("Failed to connect to database at '%s': %v", dbPath, err)
		return
	}
	defer db.Close()

	query := fmt.Sprintf(`INSERT OR REPLACE INTO "%s" (
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
		VALUES (?, ?, ?, ?,	?, ?, ?, ?,	?, ?, ?, ?, ?, ?);
	`, tableName)

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

		// Defaults for non-netscape cookie fields
		originAttr := ""     // empty string as placeholder
		inBrowserElem := 0   // always 0
		sameSite := 0        // none
		scheme := 0          // http=0 https=1
		partitionedAttr := 0 // not partitioned

		_, err = db.Exec(query,
			cookie.domain,     // host
			cookie.path,       // path
			isSecureInt,       // isSecure
			isHttpOnlyInt,     // isHttpOnly
			cookie.expiration, // expiry
			cookie.name,       // name
			cookie.value,      // value
			originAttr,        // originAttributes
			lastAccessTime,    // lastAccessed
			creationTime,      // creationTime
			inBrowserElem,     // inBrowserElement
			sameSite,          // sameSite
			scheme,            // schemeMap
			partitionedAttr,   // isPartitionedAttributeSet
		)
		if err != nil {
			err = fmt.Errorf("error updating cookie %s (%s): %v", cookie.name, cookie.domain, err)
			return
		}
	}

	return
}

/* table: cookies all columns+types
creation_utc             int
host_key                 string
top_frame_site_key       string
name                     string
value                    string
encrypted_value          []byte
path                     string
expires_utc              int
is_secure                int
is_httponly              int
last_access_utc          int
has_expires              int
is_persistent            int
priority                 int
samesite                 int
source_scheme            int
source_port              int
last_update_utc          int
source_type              int
has_cross_site_ancestor  int
*/

func writeCookiesToChrome(cookies []cookie, dbPath string, tableName string) (err error) {
	db, err := connectToDB(dbPath)
	if err != nil {
		err = fmt.Errorf("Failed to connect to database at '%s': %v", dbPath, err)
		return
	}
	defer db.Close()

	query := fmt.Sprintf(`
		INSERT OR REPLACE INTO "%s" (
			creation_utc,
			host_key,
			name,
			value,
			path,
			expires_utc,
			is_secure,
			is_httponly,
			last_access_utc,
			has_expires,
			is_persistent,
			priority,
			encrypted_value,
			samesite,
			source_scheme,
			top_frame_site_key,
			source_port,
			last_update_utc,
			source_type,
			has_cross_site_ancestor
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`, tableName)

	for _, cookie := range cookies {
		// Convert booleans
		isSecureInt := 0
		if cookie.isSecure {
			isSecureInt = 1
		}
		isHttpOnlyInt := 0
		if cookie.httpOnly {
			isHttpOnlyInt = 1
		}

		// Ensure domain has leading dot for subdomain inclusion
		if cookie.includeSub && !strings.HasPrefix(cookie.domain, ".") {
			cookie.domain = "." + cookie.domain
		}

		now := time.Now().UnixMicro()
		expiration := int64(cookie.expiration) * 1000000 // seconds to microseconds

		// Convert to chromes timestamp
		expiration = unixMicroToChromeTimestamp(expiration)
		now = unixMicroToChromeTimestamp(now)

		// Defaults for non-netscape cookie fields
		hasExpires := 1              // Does not invalidate on browser close
		isPersistent := 1            // Same as above
		priority := 1                // Medium priority
		encryptedValue := []byte("") // unused
		sameSite := -1               // unspecified
		sourceScheme := 2            // unspecified=0, http=1, https=2
		topFrameSiteKey := ""        // none
		sourcePort := 0              // unknown
		sourceType := 0              // unknown
		hasCrossSiteAncestor := 0    // not sent in cross-site iframe

		_, err = db.Exec(query,
			now,                  // creation_utc
			cookie.domain,        // host_key
			cookie.name,          // name
			cookie.value,         // value
			cookie.path,          // path
			expiration,           // expires_utc
			isSecureInt,          // is_secure
			isHttpOnlyInt,        // is_httponly
			now,                  // last_access_utc
			hasExpires,           // has_expires
			isPersistent,         // is_persistent
			priority,             // priority
			encryptedValue,       // encrypted_value
			sameSite,             // samesite
			sourceScheme,         // source_scheme
			topFrameSiteKey,      // top_frame_site_key
			sourcePort,           // source_port
			now,                  // last_update_utc
			sourceType,           // source_type
			hasCrossSiteAncestor, // has_cross_site_ancestor
		)
		if err != nil {
			err = fmt.Errorf("error updating cookie %s (%s): %v", cookie.name, cookie.domain, err)
			return
		}
	}
	return
}

func unixMicroToChromeTimestamp(unixMicro int64) int64 {
	const epochDifferenceSeconds = 11644473600
	const microsecondsPerSecond = 1_000_000

	// Convert the epoch difference to microseconds
	epochDifferenceMicro := int64(epochDifferenceSeconds * microsecondsPerSecond)

	// Add the epoch difference to the input unix microseconds
	return unixMicro + epochDifferenceMicro
}
