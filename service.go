package rincon

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type Service struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Endpoint    string    `json:"endpoint"`
	HealthCheck string    `json:"health_check"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// Service returns the current service registration of the client.
// If the client is not registered, it will be nil.
func (c *Client) Service() *Service {
	return c.service
}

// IsRegistered returns true if the client is registered.
func (c *Client) IsRegistered() bool {
	return c.service != nil
}

// Register registers the client with the given service definition and routes.
func (c *Client) Register(service Service, routes []string) (int, error) {
	req, err := c.newRequest("POST", "/rincon/services", service)
	if err != nil {
		return 0, err
	}

	newService := new(Service)
	_, apiError, err := c.do(req, newService)
	if err != nil {
		return 0, err
	} else if apiError != nil {
		return 0, fmt.Errorf("[%d] %s", apiError.StatusCode, apiError.Message)
	}

	c.service = newService
	for _, route := range routes {
		err = c.RegisterRoute(route)
		if err != nil {
			log.Printf("failed to register route %s: %s", route, err)
		}
	}
	c.StartHeartbeat()
	return newService.ID, nil
}

// Deregister de-registers the client from the Rincon server.
func (c *Client) Deregister() error {
	if c.service == nil {
		return fmt.Errorf("client is not registered")
	}

	req, err := c.newRequest("DELETE", "/rincon/services/"+strconv.Itoa(c.service.ID), nil)
	if err != nil {
		return err
	}

	_, apiError, err := c.do(req, nil)
	if err != nil {
		return err
	} else if apiError != nil {
		return fmt.Errorf("[%d] %s", apiError.StatusCode, apiError.Message)
	}

	c.service = nil
	return nil
}
