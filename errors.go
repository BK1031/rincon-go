package rincon

type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}
