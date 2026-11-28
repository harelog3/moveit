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
	// ErrTmdbUnexpectedStatus covers every non-2xx response that isn't a known
	// domain error (bad API key, rate limit, TMDB outage). The returned error
	// wraps it and includes the status.
	ErrTmdbUnexpectedStatus = errors.New("TMDB returned an unexpected status")
)

// unexpectedStatus returns an ErrTmdbUnexpectedStatus error carrying res's status.
func unexpectedStatus(res *http.Response) error {
	return fmt.Errorf("%w: %s", ErrTmdbUnexpectedStatus, res.Status)
}

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
		return nil, fmt.Errorf("parsing tmdb base url: %w", err)
	}

	u = u.JoinPath(path)
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(c, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating tmdb request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.apiKey)
	req.Header.Add("User-Agent", "moveitapi/1.0")

	return req, nil
}

func (s *Service) GetPopularMovies(c context.Context) ([]TmdbMovieSearch, error) {
	const errorContext = "getting tmdb popular movies: %w"
	req, err := s.buildHTTPRequest(c, "/movie/popular", nil)
	if err != nil {
		return nil, fmt.Errorf(errorContext, err)
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(errorContext, err)
	}

	if !httpx.IsSuccessStatus(res) {
		return nil, unexpectedStatus(res)
	}

	var result TmdbMovieSearchQuery

	err = json.UnmarshalRead(res.Body, &result)

	if err != nil {
		return nil, fmt.Errorf(errorContext, err)
	}

	return result.Results, nil
}

// SearchMoviesByTitle searches TMDB movies by title. It returns
// ErrNoTmdbSearchResults when nothing matches and an ErrTmdbUnexpectedStatus
// error for any non-2xx response.
func (s *Service) SearchMoviesByTitle(c context.Context, title string) ([]TmdbMovieSearch, error) {
	const errorContext = "searching tmdb movies: %w"

	params := url.Values{}
	params.Add("query", title)
	req, err := s.buildHTTPRequest(c, "/search/movie", params)
	if err != nil {
		return nil, fmt.Errorf(errorContext, err)
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(errorContext, err)
	}
	defer res.Body.Close()

	if !httpx.IsSuccessStatus(res) {
		return nil, fmt.Errorf(errorContext, unexpectedStatus(res))
	}
	var query TmdbMovieSearchQuery

	err = json.UnmarshalRead(res.Body, &query)

	if err != nil {
		return nil, fmt.Errorf(errorContext, err)
	}
	// TMDB answers a search with no matches with 200 and empty results.
	if len(query.Results) == 0 {
		return nil, ErrNoTmdbSearchResults
	}
	return query.Results, nil
}

// GetMovieDetailByID fetches the English details of the TMDB movie with the
// given id. It returns ErrTmdbDetailNotFound on a 404 and an
// ErrTmdbUnexpectedStatus error for any other non-2xx response.
func (s *Service) GetMovieDetailByID(c context.Context, id int) (*TmdbMovieDetail, error) {
	const errorContext = "getting tmdb movie %d: %w"
	params := url.Values{}
	params.Add("language", "en-US")
	req, err := s.buildHTTPRequest(c, "/movie/"+strconv.Itoa(id), params)
	if err != nil {
		return nil, fmt.Errorf(errorContext, id, err)
	}

	res, err := s.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf(errorContext, id, err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrTmdbDetailNotFound
	}
	if !httpx.IsSuccessStatus(res) {
		return nil, fmt.Errorf(errorContext, id, unexpectedStatus(res))
	}

	query := TmdbMovieDetail{}
	err = json.UnmarshalRead(res.Body, &query)

	if err != nil {
		return nil, fmt.Errorf(errorContext, id, err)
	}
	return &query, nil
}
