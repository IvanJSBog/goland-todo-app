package core_http_response

type ErrorResponse struct {
	Message string `json:"message" example:"short human readable message"`
	Error   string `json:"error" example:"full error text"`
}
