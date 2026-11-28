// Package httpx holds small helpers on top of net/http, mainly for
// classifying responses from external APIs such as TMDB.
package httpx

import "net/http"

// IsSuccessStatus reports whether r has a 2xx status code.
func IsSuccessStatus(r *http.Response) bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsClientError reports whether r has a 4xx status code, meaning the request
// itself is wrong and retrying it unchanged won't help.
func IsClientError(r *http.Response) bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError reports whether r has a 5xx status code, meaning the remote
// service failed and the request may succeed if retried later.
func IsServerError(r *http.Response) bool {
	return r.StatusCode >= 500 && r.StatusCode < 600
}

// IsBadRequest reports whether r is a 400, usually invalid query parameters.
func IsBadRequest(r *http.Response) bool {
	return r.StatusCode == http.StatusBadRequest
}

// IsUnauthorized reports whether r is a 401, usually a missing or invalid API
// key or token.
func IsUnauthorized(r *http.Response) bool {
	return r.StatusCode == http.StatusUnauthorized
}

// IsForbidden reports whether r is a 403: authenticated but not allowed.
func IsForbidden(r *http.Response) bool {
	return r.StatusCode == http.StatusForbidden
}

// IsNotFound reports whether r is a 404.
func IsNotFound(r *http.Response) bool {
	return r.StatusCode == http.StatusNotFound
}

// IsTooManyRequests reports whether r is a 429, i.e. the client hit a rate
// limit. The Retry-After header, when present, says how long to wait.
func IsTooManyRequests(r *http.Response) bool {
	return r.StatusCode == http.StatusTooManyRequests
}

// IsRetryable reports whether retrying the same request later may succeed:
// rate limits (429) and server errors (5xx) other than 501 Not Implemented.
func IsRetryable(r *http.Response) bool {
	if IsTooManyRequests(r) {
		return true
	}
	return IsServerError(r) && r.StatusCode != http.StatusNotImplemented
}
