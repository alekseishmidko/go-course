package core_http_response

// ErrorResponse is the common JSON error body returned by API handlers.
type ErrorResponse struct {
	Message string `json:"message" example:"failed to get task"`
	Error   string `json:"error" example:"task not found"`
}
