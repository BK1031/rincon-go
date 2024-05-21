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

func (c *Client) Rincon() *Service {
	return c.rincon
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
	c.userAgent = fmt.Sprintf("%s-%d", newService.Name, newService.ID)
	for _, route := range routes {
		err = c.RegisterRoute(route)
		if err != nil {
			log.Printf("failed to register route %s: %s", route, err)
		}
	}
	c.StartHeartbeat()
	return newService.ID, nil
}

// Deregister deregisters the client from the Rincon server.
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

func (c *Client) GetServicesByName(name string) ([]Service, error) {
	services := make([]Service, 0)
	req, err := c.newRequest("GET", "/rincon/services/"+name, services)
	if err != nil {
		return nil, err
	}

	_, apiError, err := c.do(req, &services)
	if err != nil {
		return nil, err
	} else if apiError != nil {
		return nil, fmt.Errorf("[%d] %s", apiError.StatusCode, apiError.Message)
	}

	return services, nil
}