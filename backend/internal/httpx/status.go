// Package httpx manages extension utilities for http
package httpx

import "net/http"

// IsSuccessStatus reports whether r has a 2xx status code.
func IsSuccessStatus(r *http.Response) bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}
