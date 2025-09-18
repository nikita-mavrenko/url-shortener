package utils

import "net/url"

func IsValidLink(link string) bool {
	_, err := url.Parse(link)
	if err != nil {
		return false
	}
	return true
}
