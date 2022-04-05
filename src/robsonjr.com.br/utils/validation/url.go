package validation

import (
	"errors"
	"net/url"
	"strings"
)

func getUriHostnameWithPath(parsedUrl *url.URL) string {
	if parsedUrl.RawPath != "" {
		return parsedUrl.Hostname() + parsedUrl.RawPath
	}

	return parsedUrl.Hostname() + parsedUrl.Path
}

func isValidSchema(rawUrl string) bool {
	if strings.HasPrefix(rawUrl, "http://") || strings.HasPrefix(rawUrl, "https://") {
		return true
	}

	return false
}

func isValidUrl(rawUrl string) (*url.URL, bool) {
	sanitizeUrl := rawUrl

	if !isValidSchema(rawUrl) {
		sanitizeUrl = "http://" + rawUrl
	}

	parsedUrl, err := url.ParseRequestURI(sanitizeUrl)
	if err != nil {
		return nil, false
	}

	return parsedUrl, true
}

func IsChildrenUrl(originUrl string, checkUrl string) bool {
	parsedOriginUrl, chk := isValidUrl(originUrl)
	if !chk {
		return false
	}

	parsedCheckUrl, chk := isValidUrl(checkUrl)
	if !chk {
		return false
	}

	if parsedOriginUrl.Hostname() != parsedCheckUrl.Hostname() {
		return false
	}

	originRequestUri := getUriHostnameWithPath(parsedOriginUrl)
	originCheckUri := getUriHostnameWithPath(parsedCheckUrl)

	if !strings.HasPrefix(originCheckUri, originRequestUri) {
		return false
	}

	return true
}

func GetBaseUrl(rawUrl string) (baseUrl string, err error) {
	parsedUrl, ok := isValidUrl(rawUrl);
	if !ok {
		return "", errors.New("couldn't parse url")
	}

	return parsedUrl.Scheme + "://" + parsedUrl.Host, nil
}