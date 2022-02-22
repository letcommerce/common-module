package models

import (
	"fmt"
)

type Response struct {
	Message string `json:"message" example:"Successfully added!"`
}

type StatusResponse struct {
	Status string `json:"status" example:"OK"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"Error details"`
}

func ErrorResponseF(format string, a ...interface{}) ErrorResponse {
	message := ErrorResponse{Message: fmt.Sprintf(format, a...)}
	return message
}
