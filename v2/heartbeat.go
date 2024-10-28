package rincon

import (
	"fmt"
	"log"
	"time"
)

type HeartbeatMode int32

const (
	ServerHeartbeat HeartbeatMode = 0
	ClientHeartbeat HeartbeatMode = 1
)

var heartbeatTicker *time.Ticker

// StartHeartbeat starts the heartbeat for the client.
// If the client is in server heartbeat mode, it will return an error.
// If the heartbeat is already active, it will return an error.
func (c *Client) StartHeartbeat() error {
	if c.heartbeatMode == ServerHeartbeat {
		return fmt.Errorf("client is in server heartbeat mode")
	}
	if heartbeatTicker != nil {
		return fmt.Errorf("heartbeat already active")
	}
	heartbeatTicker = c.startClientHeartbeat()
	return nil
}

// StopHeartbeat stops the heartbeat for the client.
// If the client is in server heartbeat mode, it will return an error.
// If the heartbeat is not active, it will return an error.
func (c *Client) StopHeartbeat() error {
	if c.heartbeatMode == ServerHeartbeat {
		return fmt.Errorf("client is in server heartbeat mode")
	}
	if heartbeatTicker == nil {
		return fmt.Errorf("heartbeat not active")
	}
	heartbeatTicker.Stop()
	heartbeatTicker = nil
	return nil
}

// startClientHeartbeat starts the client heartbeat.
// It returns a ticker that can be used to stop the heartbeat.
func (c *Client) startClientHeartbeat() *time.Ticker {
	ticker := time.NewTicker(time.Duration(c.heartbeatInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				id, err := c.Register(*c.service, []Route{})
				if err != nil {
					log.Printf("heartbeat failed: %s", err)
				} else {
					log.Printf("heartbeat success: %d", id)
				}
			}
		}
	}()
	return ticker
}
