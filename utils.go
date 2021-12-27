package forward

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	xurls "mvdan.cc/xurls/v2"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

var (
	urlRegexp *regexp.Regexp
)

func init() {
	if u, err := xurls.StrictMatchingScheme("https?://|wss?://"); err != nil {
		panic(err)
	} else {
		urlRegexp = u
	}
}

func isHttpUrl(u string) bool {
	return regexp.MustCompile(`^https?:\/\/`).MatchString(u)
}

func replaceHost(content, oldHost, newHost string) string {
	if !strings.HasPrefix(oldHost, "http") {
		oldHost = "http://" + oldHost
	}
	if !strings.HasPrefix(newHost, "http") {
		newHost = "http://" + newHost
	}

	newContent := content

	newHostUrl, err := url.Parse(newHost)

	if err != nil {
		panic(err)
	}

	oldHostUrl, err := url.Parse(oldHost)

	if err != nil {
		panic(err)
	}

	newContent = urlRegexp.ReplaceAllStringFunc(newContent, func(s string) string {
		matchUrl, err := url.Parse(s)

		if err != nil {
			return s
		}

		query := []string{}
		queryArr := strings.Split(matchUrl.RawQuery, "&")

		for _, q := range queryArr {
			arr := strings.Split(q, "=")
			key := arr[0]
			if len(arr) == 1 {
				if strings.Contains(q, "=") {
					query = append(query, key+"=")
				} else {
					query = append(query, key)
				}
			} else {
				escapedValue := strings.Join(arr[1:], "=")

				if unescapedValue, err := url.QueryUnescape(escapedValue); err == nil {
					escapedValue = url.QueryEscape(replaceHost(unescapedValue, oldHost, newHost))
				} else {
					escapedValue = replaceHost(escapedValue, oldHost, newHost)
				}

				query = append(query, key+"="+escapedValue)
			}
		}

		matchUrl.RawQuery = strings.Join(query, "&")

		if matchUrl.Host != oldHostUrl.Host {
			if contains([]string{"http", "https"}, matchUrl.Scheme) || strings.HasPrefix(s, "//") {
				return fmt.Sprintf("%s://%s/?forward_url=%s", newHostUrl.Scheme, newHostUrl.Host, url.QueryEscape(matchUrl.String()))
			} else if contains([]string{"ws", "wss"}, matchUrl.Scheme) {
				return fmt.Sprintf("%s://%s/?forward_url=%s", "ws", newHostUrl.Host, url.QueryEscape(matchUrl.String()))
			}

			return s
		}

		if matchUrl.Scheme == "https" {
			matchUrl.Scheme = "http"
		} else if matchUrl.Scheme == "wss" {
			matchUrl.Scheme = "ws"
		}

		matchUrl.Host = newHostUrl.Host

		return matchUrl.String()
	})

	newContent = strings.ReplaceAll(newContent, "//"+oldHostUrl.Host, "//"+newHostUrl.Host)

	return newContent
}
