package request

import (
	"errors"
	"net/http"
	"strings"
)

const BearerPrefix string = "Bearer "

var (
	ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")
	ErrMalformedAuthHeader  = errors.New("malformed authorization header")
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return "", ErrMalformedAuthHeader
	}

	token := authHeader[len(BearerPrefix):]
	if token == "" {
		return "", ErrMalformedAuthHeader
	}

	return token, nil
}
