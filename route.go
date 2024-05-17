package rincon

import (
	"fmt"
	"time"
)

type Route struct {
	Route       string    `json:"route"`
	ServiceName string    `json:"service_name"`
	CreatedAt   time.Time `json:"created_at"`
}

func (c *Client) Routes() ([]Route, error) {
	if c.service == nil {
		return nil, fmt.Errorf("client is not registered")
	}
	return c.RoutesForService(c.service.Name)
}

func (c *Client) RegisterRoute(route string) error {
	if c.service == nil {
		return fmt.Errorf("client is not registered")
	}
	req, err := c.newRequest("POST", "/rincon/routes", Route{
		Route:       route,
		ServiceName: c.service.Name,
	})
	if err != nil {
		return err
	}

	_, apiError, err := c.do(req, nil)
	if err != nil {
		return err
	} else if apiError != nil {
		return fmt.Errorf("[%d] %s", apiError.StatusCode, apiError.Message)
	}

	return nil
}

func (c *Client) RoutesForService(serviceName string) ([]Route, error) {
	req, err := c.newRequest("GET", "/rincon/services/"+serviceName+"/routes", nil)
	if err != nil {
		return nil, err
	}

	var routes []Route
	_, _, err = c.do(req, &routes)
	if err != nil {
		return nil, err
	}

	return routes, nil
}
