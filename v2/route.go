package rincon

import (
	"fmt"
	"strings"
	"time"
)

type Route struct {
	Route       string    `json:"route"`
	ServiceName string    `json:"service_name"`
	Method      string    `json:"method"`
	CreatedAt   time.Time `json:"created_at"`
}

// Routes returns the routes registered for the client's service.
// If the client is not registered, an error will be returned.
func (c *Client) Routes() ([]Route, error) {
	if c.service == nil {
		return nil, fmt.Errorf("client is not registered")
	}
	return c.RoutesForService(c.service.Name)
}

// RegisterRoute registers a route for the client's service.
// If the client is not already registered, an error will be returned.
func (c *Client) RegisterRoute(route string, method string) error {
	if c.service == nil {
		return fmt.Errorf("client is not registered")
	}
	req, err := c.newRequest("POST", "/rincon/routes", Route{
		Route:       route,
		ServiceName: c.service.Name,
		Method:      method,
	}, nil)
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

// MatchRoute returns the service that is registered to handle the given route.
func (c *Client) MatchRoute(route string, method string) (*Service, error) {
	if c.service == nil {
		return nil, fmt.Errorf("client is not registered")
	}
	route = strings.TrimPrefix(route, "/")
	route = strings.TrimSuffix(route, "/")
	req, err := c.newRequest("GET", "/rincon/match", nil, map[string]string{
		"route":  route,
		"method": method,
	})
	if err != nil {
		return nil, err
	}

	var service Service
	_, apiError, err := c.do(req, &service)
	if err != nil {
		return nil, err
	} else if apiError != nil {
		return nil, fmt.Errorf("[%d] %s", apiError.StatusCode, apiError.Message)
	}
	return &service, nil
}

// RoutesForService returns the routes registered for the given service.
func (c *Client) RoutesForService(serviceName string) ([]Route, error) {
	req, err := c.newRequest("GET", "/rincon/services/"+serviceName+"/routes", nil, nil)
	if err != nil {
		return nil, err
	}

	var routes []Route
	_, apiError, err := c.do(req, &routes)
	if err != nil {
		return nil, err
	} else if apiError != nil {
		return nil, fmt.Errorf("[%d] %s", apiError.StatusCode, apiError.Message)
	}

	return routes, nil
}
