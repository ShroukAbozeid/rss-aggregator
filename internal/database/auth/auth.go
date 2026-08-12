package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Example:
// X-API-Key: ApiKey {insert api key here}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("X-API-Key")
	if val == "" {
		return "", errors.New("Missing API key")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("Invalid API key")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("Invalid API key")
	}

	return vals[1], nil
}
