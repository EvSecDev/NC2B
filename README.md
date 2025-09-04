# Netscape Cookie to Browser (NC2B)

## Overview

A simple CLI utility to take cookies from a netscape cookie text file and write them into browser cookie store files.

Currently only supports writing to Firefox cookies.sqlite file.

## NC2B Help Menu

```bash
Netscape Cookie to Browser (NC2B)
  Retrieves cookies from netscape text files and injects into browser cookie storage

  Options:
    -j, --cookie-jar         <path/to/file>  Path to netscape cookie file to read cookies
    -b, --browser-cookie-jar <path/to/file>  Path to browser cookie storage file
    -n, --name               <text>          Import cookies matching name
    -d, --domain             <text>          Import cookies matching hostname
        --allow-expired                      Enable importing expired cookies
    -h, --help                               Show this help menu
    -V, --version                            Show version and packages
        --versionid                          Show only version number

Report bugs to: dev@evsec.net
NC2B home page: <https://github.com/EvSecDev/NC2B>
General help using GNU software: <https://www.gnu.org/gethelp/>
```

## Notes

### Filtering

Cookie name is partial match, meaning if the provided filter text is present in the cookie name it matches.

Cookie domain is suffix match only, meaning if the provided filter text is present, anchored to the right, it matches.

### Browser Cookie Limitations

Due to netscape cookie files lacking the additional fields that many browsers natively have in their own cookie store, this program sets defaults for certain fields.

Below are the fields and their default values.

Firefox:

- `originAttributes` = ''
- `inBrowserElement` = `0`
- `sameSite` = `0`
- `schemeMap` = `0`
- `isPartitionedAttributeSet` = `0`
