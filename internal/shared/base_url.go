package shared

import (
	"strings"
)

func removeFirstChar(s string) string {
	if len(s) > 0 {
		return s[1:]
	}
	return s
}

func ConcatBaseUrlAndUri(baseUrl string, uri string) string {
	if baseUrl == "" {
		panic("make sure that you update the environment variables")
	}

	if strings.HasSuffix(baseUrl, "/") {
		uri = removeFirstChar(uri)
	}
	return baseUrl + uri
}
