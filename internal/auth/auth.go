package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey, ok := strings.CutPrefix(headers.Get("Authorization"), "ApiKey ")
	if !ok {
		return "", errors.New("No API key provided")
	}

	return apiKey, nil
}
