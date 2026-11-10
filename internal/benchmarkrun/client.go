package benchmarkrun

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// requestTimeout limits each request of the client, except a snapshot download, which can take minutes
const requestTimeout = 30 * time.Second

// Client runs benchmarks on a powerpipe server, through its benchmark run API
type Client struct {
	host *url.URL
	http *http.Client
}

// APIError is an error response from the server
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Message
}

// NewClient returns a client for the powerpipe server at host, e.g. http://localhost:9033
func NewClient(host string) (*Client, error) {
	u, err := url.Parse(host)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, fmt.Errorf("invalid host '%s': expected a URL such as http://localhost:9033", host)
	}
	return &Client{host: u, http: &http.Client{}}, nil
}

// Start starts a run of the benchmark, and returns the run and its URL
func (c *Client) Start(ctx context.Context, benchmark string, timeout time.Duration) (Run, *url.URL, error) {
	body, err := json.Marshal(map[string]int{"timeout_seconds": int(timeout / time.Second)})
	if err != nil {
		return Run{}, nil, err
	}
	startURL := c.host.JoinPath("api", "v0", "benchmark", benchmark, "run")
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	resp, err := c.do(ctx, http.MethodPost, startURL, bytes.NewReader(body), http.StatusAccepted)
	if err != nil {
		return Run{}, nil, err
	}
	defer resp.Body.Close()
	var run Run
	if err := json.NewDecoder(resp.Body).Decode(&run); err != nil {
		return Run{}, nil, fmt.Errorf("invalid response from %s: %w", startURL.Redacted(), err)
	}
	runURL, err := startURL.Parse(resp.Header.Get("Location"))
	if err != nil || resp.Header.Get("Location") == "" {
		return Run{}, nil, fmt.Errorf("the response from %s has no valid Location header", startURL.Redacted())
	}
	return run, runURL, nil
}

// Get returns the state of the run
func (c *Client) Get(ctx context.Context, runURL *url.URL) (Run, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	resp, err := c.do(ctx, http.MethodGet, runURL, nil, http.StatusOK)
	if err != nil {
		return Run{}, err
	}
	defer resp.Body.Close()
	var run Run
	if err := json.NewDecoder(resp.Body).Decode(&run); err != nil {
		return Run{}, fmt.Errorf("invalid response from %s: %w", runURL.Redacted(), err)
	}
	return run, nil
}

// Snapshot returns the snapshot of a completed run. The caller must close it.
func (c *Client) Snapshot(ctx context.Context, runURL *url.URL) (io.ReadCloser, error) {
	resp, err := c.do(ctx, http.MethodGet, runURL.JoinPath("snapshot"), nil, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// Delete cancels the run if it is running, and removes it from the server
func (c *Client) Delete(ctx context.Context, runURL *url.URL) error {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	resp, err := c.do(ctx, http.MethodDelete, runURL, nil, http.StatusNoContent)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

// do sends the request, and returns an error for a status other than want
func (c *Client) do(ctx context.Context, method string, u *url.URL, body io.Reader, want int) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		// the url.Error repeats the method and the URL of the request
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			err = urlErr.Err
		}
		return nil, fmt.Errorf("could not reach the powerpipe server at %s: %w", c.host.Redacted(), err)
	}
	if resp.StatusCode == want {
		return resp, nil
	}
	defer resp.Body.Close()
	// the server sends errors as {"title": ..., "detail": ...}
	var errorBody struct {
		Title  string `json:"title"`
		Detail string `json:"detail"`
	}
	message := resp.Status
	if data, readErr := io.ReadAll(io.LimitReader(resp.Body, 64*1024)); readErr == nil && json.Unmarshal(data, &errorBody) == nil {
		if errorBody.Detail != "" {
			message = errorBody.Detail
		} else if errorBody.Title != "" {
			message = errorBody.Title
		}
	}
	return nil, &APIError{StatusCode: resp.StatusCode, Message: message}
}
