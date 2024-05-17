package rincon

// ErrorResponse is a struct to help decode errors from the Rincon API.
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}
