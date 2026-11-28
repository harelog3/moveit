package tmdb

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/harelog3/moveit/internal/httpx"
)

var (
	ErrNoTmdbSearchResults = errors.New("TMDB returned 0 results for the given title")
	ErrTmdbDetailNotFound  = errors.New("TMDB with given id was not found")
)

// api is the TMDB v3 base URL. Requests authenticate with a v4 read access
// token sent as a Bearer header, not with the api_key query parameter.
const api = "https://api.tmdb.org/3"

// Service is a TMDB API client. It is safe for concurrent use.
type Service struct {
	apiKey string
	client *http.Client
}

// NewService returns a TMDB client. The transport and client timeouts bound
// every request, so a slow TMDB cannot hang the handlers that call it.
func NewService(apiKey string) *Service {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,  // Time to establish TCP connection
			KeepAlive: 30 * time.Second, // Keep-alive probe interval
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,  // Time for TLS handshake
		ResponseHeaderTimeout: 10 * time.Second, // Time to wait for response headers
		MaxIdleConns:          100,              // Total idle connections across all hosts
		MaxIdleConnsPerHost:   10,               // Idle connections per single host (default is 2)
		IdleConnTimeout:       90 * time.Second, // How long an idle connection stays open
	}
	return &Service{apiKey: apiKey, client: &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}}
}

// buildHTTPRequest builds a GET request to path under the TMDB base URL with
// the auth and accept headers TMDB expects.
func (s *Service) buildHTTPRequest(c context.Context, path string, params url.Values) (*http.Request, error) {
	u, err := url.Parse(api)
	if err != nil {
		return nil, err
	}

	u = u.JoinPath(path)
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(c, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.apiKey)
	req.Header.Add("User-Agent", "moveitapi/1.0")

	return req, nil
}

// SearchMoviesByTitle searches TMDB movies by title. Any non-2xx response is
// reported as ErrNoTmdbSearchResults.
func (s *Service) SearchMoviesByTitle(c context.Context, title string) ([]TmdbMovieSearch, error) {
	params := url.Values{}
	params.Add("query", title)
	req, err := s.buildHTTPRequest(c, "/search/movie", params)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if !httpx.IsSuccessStatus(res) {
		fmt.Println(res.Status)
		return nil, ErrNoTmdbSearchResults
	}
	var query TmdbMovieSearchQuery

	err = json.UnmarshalRead(res.Body, &query)

	if err != nil {
		return nil, err
	}
	return query.Results, nil
}

// GetMovieDetailByID fetches the English details of the TMDB movie with the
// given id. Any non-2xx response is reported as ErrTmdbDetailNotFound.
func (s *Service) GetMovieDetailByID(c context.Context, id int) (*TmdbMovieDetail, error) {
	params := url.Values{}
	params.Add("language", "en-US")
	req, err := s.buildHTTPRequest(c, "/movie/"+strconv.Itoa(id), params)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if !httpx.IsSuccessStatus(res) {
		return nil, ErrTmdbDetailNotFound
	}

	query := TmdbMovieDetail{}
	err = json.UnmarshalRead(res.Body, &query)

	if err != nil {
		return nil, err
	}
	return &query, nil
}
