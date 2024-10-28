package rincon

import "fmt"

type Ping struct {
	Message  string `json:"message"`
	Routes   int    `json:"routes"`
	Services int    `json:"services"`
}

// Ping sends a request to the Rincon server to ensure its reachable.
func (c *Client) Ping() (*Ping, error) {
	req, err := c.newRequest("GET", "/rincon/ping", nil, nil)
	if err != nil {
		return nil, err
	}

	ping := new(Ping)
	_, apiError, err := c.do(req, ping)
	if err != nil {
		return nil, err
	} else if apiError != nil {
		return nil, fmt.Errorf("[%d] %s", apiError.StatusCode, apiError.Message)
	}

	return ping, nil
}
