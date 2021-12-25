package forward

import (
	"net/url"
	"regexp"
	"strings"
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
	urlRegexp = regexp.MustCompile(`(((http|ws)s?:)?//)?(([\d\w]|%[a-fA-f\d]{2,2})+(:([\d\w]|%[a-fA-f\d]{2,2})+)?@)?([\d\w][-\d\w]{0,253}[\d\w]\.)+[\w]{2,63}(:[\d]+)?(/([-+_~.\d\w]|%[a-fA-f\d]{2,2})*)*(\?(&?([-+_~.\d\w]|%[a-fA-f\d]{2,2})=?)*)?(#([-+_~.\d\w]|%[a-fA-f\d]{2,2})*)?`)
)

func replaceHost(content, origin, target string) string {
	if !strings.HasPrefix(origin, "http") {
		origin = "http://" + origin
	}
	if !strings.HasPrefix(target, "http") {
		target = "http://" + target
	}

	newContent := content

	t, err := url.Parse(target)

	if err != nil {
		panic(err)
	}

	o, err := url.Parse(origin)

	if err != nil {
		panic(err)
	}

	newContent = urlRegexp.ReplaceAllStringFunc(newContent, func(s string) string {
		u, err := url.Parse(s)

		if err != nil {
			return s
		}

		query := []string{}
		queryArr := strings.Split(u.RawQuery, "&")

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
					escapedValue = url.QueryEscape(replaceHost(unescapedValue, origin, target))
				} else {
					escapedValue = replaceHost(escapedValue, origin, target)
				}

				query = append(query, key+"="+escapedValue)
			}
		}

		u.RawQuery = strings.Join(query, "&")

		if u.Host != o.Host {
			return u.String()
		}

		if u.Scheme == "https" {
			u.Scheme = "http"
		} else if u.Scheme == "wss" {
			u.Scheme = "ws"
		}

		u.Host = t.Host

		return u.String()
	})

	return newContent
}
