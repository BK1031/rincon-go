package rincon

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL           *url.URL
	heartbeatMode     HeartbeatMode
	heartbeatInterval int
	authUser          string
	authPassword      string
	userAgent         string
	httpClient        *http.Client
}

type Config struct {
	BaseURL           string
	HeartbeatMode     HeartbeatMode
	HeartbeatInterval int
	AuthUser          string
	AuthPassword      string
}

func New(config Config) (*Client, error) {
	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, err
	}
	client := &Client{
		baseURL:           baseURL,
		heartbeatMode:     config.HeartbeatMode,
		heartbeatInterval: config.HeartbeatInterval,
		authUser:          config.AuthUser,
		authPassword:      config.AuthPassword,
		userAgent:         "rincon-go",
		httpClient:        &http.Client{},
	}
	if _, err = client.Ping(); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.baseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	return req, nil
}
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
