package response

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Response struct {
	Message string `json:"message" example:"Successfully added!"`
	ID      uint   `json:"id,omitempty" example:"1"` // the effected ID (if exists)
}

type StatusResponse struct {
	Status string `json:"status" example:"OK"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"Error details"`
	Error   string `json:"error,omitempty"`
}

func ErrorResponseF(format string, a ...interface{}) ErrorResponse {
	message := ErrorResponse{Message: fmt.Sprintf(format, a...)}
	return message
}

func NewErrorResponseF(err error, format string, a ...interface{}) ErrorResponse {
	message := fmt.Sprintf(format, a...)
	tryGetStackTrace(message, err)
	return ErrorResponse{Message: message, Error: err.Error()}
}

func NewErrorResponse(message string, err error) ErrorResponse {
	tryGetStackTrace(message, err)
	return ErrorResponse{Message: message, Error: err.Error()}
}

func tryGetStackTrace(message string, err error) string {
	p, _ := os.Getwd()
	var errorDetails string
	stackTrace := strings.ReplaceAll(fmt.Sprintf("%+v", err), p, "") // removing the working directory to make it more readable
	if stackTrace != "" {
		prefix := strings.Index(stackTrace, "\n") + 1
		if prefix < len(stackTrace) {
			last := stackTrace[prefix:]
			last = strings.ReplaceAll(last, "\n", "  ")
			split := strings.Split(last, "\t")
			if len(split) > 2 { // taking only the first two levels of the stacktrace (if exists)
				errorDetails = strings.Join(split[0:2], "  ")
				errorDetails = strings.TrimSuffix(errorDetails, "  ")
			}
		}
	}
	if errorDetails != "" {
		log.Errorf("Returning error with message: %v, error: %v, stackTrace: %v", message, err.Error(), stackTrace)
	}
	return errorDetails
}
