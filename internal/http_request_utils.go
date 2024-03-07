package internal

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func getQueryParam(name string, r *http.Request) (string, bool) {
	values, ok := r.URL.Query()[name]

	if !ok || len(values) < 1 {
		return "", false
	}

	return values[0], true
}

func getQueryParamWithDefault(name string, defaultValue string, r *http.Request) string {
	value, ok := getQueryParam(name, r)

	if !ok {
		return defaultValue
	}

	return value
}

func getQueryParamInt(name string, r *http.Request) (int, bool) {
	value, ok := getQueryParam(name, r)

	if !ok {
		return 0, false
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}

	return intValue, true
}

func getQueryParamIntWithDefault(name string, defaultValue int, r *http.Request) int {
	value, ok := getQueryParamInt(name, r)

	if !ok {
		return defaultValue
	}

	return value
}

func getQueryParamBool(name string, r *http.Request) (bool, bool) {
	value, ok := getQueryParam(name, r)

	if !ok {
		return false, false
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, false
	}

	return boolValue, true
}

func getQueryParamBoolWithDefault(name string, defaultValue bool, r *http.Request) bool {
	value, ok := getQueryParamBool(name, r)

	if !ok {
		return defaultValue
	}

	return value
}

func removeParamFromQuery(r *http.Request, paramName string) {
	queryValues := r.URL.Query()
	queryValues.Del(paramName)
	r.URL.RawQuery = queryValues.Encode()
}

func currentRequestUrl(r *http.Request) string {
	return currentRequestHost(r) + r.URL.String()
}

func currentRequestHost(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}

func escapeURLPath(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}

func ensureValidPathPrefixFormat(path string) string {
	if path == "" {
		return ""
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return strings.TrimSuffix(path, "/")
}
