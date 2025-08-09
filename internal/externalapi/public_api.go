package api

import (
	"net/url"
)

func BuildURL(baseURI string, params map[string]string) (string, error) {
	parsedURL, err := url.Parse(baseURI)
	if err != nil {
		return "", err
	}

	query := parsedURL.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}
