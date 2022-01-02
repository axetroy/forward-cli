package forward

import (
	"fmt"
	"net"
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
	urlWithSchemeRegExp    *regexp.Regexp
	hostNameRegexp         = regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})$`)
	rewriteContentExtNames = map[string]struct{}{
		".html":  {},
		".htm":   {},
		".xhtml": {},
		".xml":   {},
		".yml":   {},
		".yaml":  {},
		".css":   {},
		".js":    {},
		".txt":   {},
		".text":  {},
		".json":  {},
	}
	htmlExtNames = map[string]struct{}{
		".html":  {},
		".htm":   {},
		".xhtml": {},
	}
)

func init() {
	if u, err := xurls.StrictMatchingScheme("https?://|wss?://|//"); err != nil {
		panic(err)
	} else {
		urlWithSchemeRegExp = u
	}
}

func isHttpUrl(u string) bool {
	return regexp.MustCompile(`^https?:\/\/`).MatchString(u)
}

func isShouldReplaceContent(extNames []string) bool {
	for _, extName := range extNames {
		if _, ok := rewriteContentExtNames[extName]; ok {
			return true
		}
	}
	return false
}

func isHtml(extNames []string) bool {
	for _, extName := range extNames {
		if _, ok := htmlExtNames[extName]; ok {
			return true
		}
	}
	return false
}

func replaceHost(content, oldHost, newHost string, useSSL bool, proxyExternal bool, proxyExternalIgnores []string) string {
	newContent := urlWithSchemeRegExp.ReplaceAllStringFunc(content, func(s string) string {
		matchUrl, err := url.Parse(s)

		if err != nil {
			return s
		}

		// if host not a IP address or a valid domain name
		if net.ParseIP(matchUrl.Hostname()) == nil && !hostNameRegexp.MatchString(matchUrl.Hostname()) {
			return s
		}

		// overide url in query
		{
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
						escapedValue = url.QueryEscape(replaceHost(unescapedValue, oldHost, newHost, useSSL, proxyExternal, proxyExternalIgnores))
					} else {
						escapedValue = replaceHost(escapedValue, oldHost, newHost, useSSL, proxyExternal, proxyExternalIgnores)
					}

					query = append(query, key+"="+escapedValue)
				}
			}

			matchUrl.RawQuery = strings.Join(query, "&")
		}

		// if the host not match the target
		if matchUrl.Host != oldHost {
			// do not proxy external link
			if !proxyExternal {
				return s
			}

			// ignore proxy for this domain
			if contains(proxyExternalIgnores, matchUrl.Host) {
				return s
			}

			if contains([]string{"http", "https"}, matchUrl.Scheme) {
				scheme := "http"
				if useSSL {
					scheme = "https"
				}
				return fmt.Sprintf("%s://%s/?forward_url=%s", scheme, newHost, url.QueryEscape(matchUrl.String()))
			} else if contains([]string{"ws", "wss"}, matchUrl.Scheme) {
				scheme := "ws"
				if useSSL {
					scheme = "wss"
				}
				return fmt.Sprintf("%s://%s/?forward_url=%s", scheme, newHost, url.QueryEscape(matchUrl.String()))
			}

			return s
		}

		if contains([]string{"http", "https"}, matchUrl.Scheme) {
			if useSSL {
				s = regexp.MustCompile(`^https?:\/\/`).ReplaceAllString(s, "https://")
			} else {
				s = regexp.MustCompile(`^https?:\/\/`).ReplaceAllString(s, "http://")
			}
		} else if contains([]string{"ws", "wss"}, matchUrl.Scheme) {
			if useSSL {
				s = regexp.MustCompile(`^wss?:\/\/`).ReplaceAllString(s, "wss://")
			} else {
				s = regexp.MustCompile(`^wss?:\/\/`).ReplaceAllString(s, "ws://")
			}
		}

		s = strings.Replace(s, oldHost, newHost, 1)

		return s
	})

	return newContent
}
