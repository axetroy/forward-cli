package forward

import (
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

func replaceHost(origin, target string) string {
	newContent := origin

	replaceMap := map[string]string{
		"https://": "http://",
		"wss://":   "ws://",
		"//":       "//",
		"":         "",
	}

	for oldScheme, newScheme := range replaceMap {
		originHost := oldScheme + origin
		proxyHost := newScheme + target

		// newBodyStr = strings.ReplaceAll(newBodyStr, originHost, proxyHost)
		newContent = regexp.MustCompile("\\b"+strings.ReplaceAll(originHost, ".", "\\.")+"\\b").ReplaceAllString(newContent, proxyHost)
	}

	return newContent
}
