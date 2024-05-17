package rincon

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Client represents a client to the Rincon API.
type Client struct {
	baseURL           *url.URL
	heartbeatMode     HeartbeatMode
	heartbeatInterval int
	authUser          string
	authPassword      string
	userAgent         string
	httpClient        *http.Client
	service           *Service
}

// Config represents the configuration for a Rincon Client.
// BaseURL is the URL of the Rincon server.
// HeartbeatMode is the mode of the heartbeat. If it is set to
// ServerHeartbeat, the Rincon server will send heartbeats to the client.
// If it is set to ClientHeartbeat, the client will send heartbeats to the server.
// HeartbeatInterval is the interval of the heartbeat.
// AuthUser is the username for authentication.
// AuthPassword is the password for authentication.
type Config struct {
	BaseURL           string
	HeartbeatMode     HeartbeatMode
	HeartbeatInterval int
	AuthUser          string
	AuthPassword      string
}

// NewClient creates a new Rincon Client with the given Config.
// It returns an error if the BaseURL is invalid or the
// client cannot connect to the Rincon server.
func NewClient(config Config) (*Client, error) {
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
	req.SetBasicAuth(c.authUser, c.authPassword)
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, *ErrorResponse, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		if v != nil {
			err = json.Unmarshal(body, v)
		}
	} else {
		respError := new(ErrorResponse)
		err = json.Unmarshal(body, respError)
		respError.StatusCode = resp.StatusCode
		return resp, respError, err
	}
	return resp, nil, err
}
