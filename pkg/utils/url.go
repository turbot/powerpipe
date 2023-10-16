package utils

import (
	"strings"
)

// Modifies(trims) the URL if contains http ot https in arguments
func TrimGitUrls(urls []string) []string {
	res := make([]string, len(urls))
	for i, url := range urls {
		res[i] = strings.TrimPrefix(url, "http://")
		res[i] = strings.TrimPrefix(url, "https://")
	}
	return res
}
