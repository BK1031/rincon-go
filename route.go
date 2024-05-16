package rincon

import "time"

type Route struct {
	Route       string    `json:"route"`
	ServiceName string    `json:"service_name"`
	CreatedAt   time.Time `json:"created_at"`
}
