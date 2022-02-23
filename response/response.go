package response

import (
	"fmt"
)

type Response struct {
	Message string `json:"message" example:"Successfully added!"`
	ID      uint   `json:"id;omitempty" example:"1"` // the effected ID
}

type StatusResponse struct {
	Status string `json:"status" example:"OK"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"Error details"`
	Error   string `json:"error;omitempty"`
}

func ErrorResponseF(format string, a ...interface{}) ErrorResponse {
	message := ErrorResponse{Message: fmt.Sprintf(format, a...)}
	return message
}
