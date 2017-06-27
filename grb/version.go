package grb

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	"github.com/pkg/errors"
)

type GoVersion struct {
	Major   int
	Minor   int
	Patch   int
	Release string
	URL     string
}

func GetVersions() (stable GoVersion, unstable GoVersion, err error) {
	// Get all versions from the download page
	versions, err := getAllVersions()
	if err != nil {
		return
	}

	// Find the latest stable and unstable versions
	for _, v := range versions {
		stbl := isStable(v)
		if stbl && compare(v, stable) {
			stable = v
			continue
		}
		if !stbl && compare(v, unstable) {
			unstable = v
		}
	}

	return
}

func getAllVersions() ([]GoVersion, error) {
	resp, err := http.Get("https://golang.org/dl/")
	if err != nil {
		return nil, errors.Wrap(err, "go download page check failed")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("go download page unexpected status code")
	}

	// Get download links
	versions := make([]GoVersion, 0)
	z := html.NewTokenizer(resp.Body)
OUT:
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			break OUT
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if anchor tag
			if t.Data == "a" {
				// Search the attributes for the href
				for _, a := range t.Attr {
					if a.Key == "href" {
						// If this is a src.tar.gz file, parse it and save a struct
						if vers, ok := parseVersion(a.Val); ok {
							versions = append(versions, vers)
						}
						break
					}
				}
			}
		}
	}

	return versions, nil
}

func parseVersion(url string) (vers GoVersion, ok bool) {
	if !strings.HasSuffix(url, ".src.tar.gz") {
		return GoVersion{}, false
	}

	r := regexp.MustCompile(`go(\d)\.(\d)((\.|rc|beta)(\d))?`)
	matches := r.FindAllStringSubmatch(url, -1)

	for _, match := range matches {
		//     0         1   2    3      4      5
		// |go1.3.1|    |1| |3| |.1|    |.|    |1|
		// |go1.3|      |1| |3| ||      ||     ||
		// |go1.2.2|    |1| |2| |.2|    |.|    |2|
		// |go1.9beta2| |1| |9| |beta2| |beta| |2|

		major, _ := strconv.Atoi(match[1])
		minor, _ := strconv.Atoi(match[2])
		patch, _ := strconv.Atoi(match[5])

		vers = GoVersion{
			Major: major,
			Minor: minor,
			Patch: patch,
			URL:   url,
		}

		if match[4] != "" && match[4] != "." {
			vers.Release = match[4]
		}

		ok = true
	}

	return
}

// Return true if version a is greater than version b
func compare(a GoVersion, b GoVersion) bool {
	rel := []string{"beta", "rc", ""}

	if a.Major > b.Major {
		return true
	} else if a.Major < b.Major {
		return false
	}

	if a.Minor > b.Minor {
		return true
	} else if a.Minor < b.Minor {
		return false
	}

	// Compare releases: "" > rc > beta
	var (
		ar int
		br int
	)
	for i, v := range rel {
		if a.Release == v {
			ar = i + 1
		}
		if b.Release == v {
			br = i + 1
		}
	}
	if ar > br {
		return true
	} else if ar < br {
		return false
	}

	// We compare the patch number last because beta2 is saved as patch = 2, release = beta
	if a.Patch > b.Patch {
		return true
	} else if a.Patch < b.Patch {
		return false
	}

	return false
}

func isStable(vers GoVersion) bool {
	return vers.Release == ""
}
