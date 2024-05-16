package rincon

type Ping struct {
	Message  string `json:"message"`
	Routes   int    `json:"routes"`
	Services int    `json:"services"`
}

func (c *Client) Ping() (*Ping, error) {
	req, err := c.newRequest("GET", "/ping", nil)
	if err != nil {
		return nil, err
	}

	ping := new(Ping)
	_, err = c.do(req, ping)
	if err != nil {
		return nil, err
	}

	return ping, nil
}
