package todi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Client wraps Todoist API calls.
type Client struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
	Verbose bool
}

// NewClient creates a Todoist API client.
func NewClient(baseURL, token string, verbose bool) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Token:   token,
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
		Verbose: verbose,
	}
}

type listResponse[T any] struct {
	Results    []T    `json:"results"`
	NextCursor string `json:"next_cursor"`
}

func (c *Client) get(ctx context.Context, path string, params map[string]string, out any) error {
	fullURL, err := c.url(path, params)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func (c *Client) post(ctx context.Context, path string, body map[string]any, out any) ([]byte, error) {
	fullURL, err := c.url(path, nil)
	if err != nil {
		return nil, err
	}
	var payload io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		payload = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, payload)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.do(req, out)
}

func (c *Client) delete(ctx context.Context, path string) ([]byte, error) {
	fullURL, err := c.url(path, nil)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fullURL, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req, nil)
}

func (c *Client) deleteWithParams(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	fullURL, err := c.url(path, params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fullURL, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req, nil)
}

func (c *Client) url(path string, params map[string]string) (string, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", err
	}
	u.Path = strings.TrimRight(u.Path, "/") + path
	if len(params) > 0 {
		q := u.Query()
		for key, value := range params {
			if value == "" {
				continue
			}
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
	}
	return u.String(), nil
}

func (c *Client) do(req *http.Request, out any) ([]byte, error) {
	req.Header.Set("Authorization", "Bearer "+c.Token)
	if c.Verbose {
		writef(os.Stderr, "%s %s\n", req.Method, req.URL.String())
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("api error: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	if out != nil && len(body) > 0 {
		if err := json.Unmarshal(body, out); err != nil {
			return body, fmt.Errorf("decode response: %w", err)
		}
	}
	return body, nil
}
