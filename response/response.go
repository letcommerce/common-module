package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	requestid "github.com/sumit-tembe/gin-requestid"
	"os"
	"strings"
)

var (
	ctx *gin.Context
)

func Init(ginCtx *gin.Context) {
	ctx = ginCtx
}

type Response struct {
	Message string `json:"message" example:""`
	ID      uint   `json:"id,omitempty" example:"1"` // the effected ID (if exists)
}

type StatusResponse struct {
	Status string `json:"status" example:"OK"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"Error details"`
	Error   string `json:"error,omitempty"`
}

func NewResponse(message string, id uint) Response {
	return Response{Message: message}
}

func ErrorResponseF(format string, a ...interface{}) ErrorResponse {
	return ErrorResponse{Message: fmt.Sprintf(format, a...)}
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
		split := strings.Split(stackTrace, "\n")
		if len(split) > 6 {
			errorDetails = strings.Join(split[0:6], "  \n")
		} else {
			errorDetails = stackTrace
		}
	}
	if errorDetails != "" {
		log.Errorf("Returning error with message: %v, error: %v, stackTrace: %v [%v]", message, err.Error(), errorDetails, requestid.GetRequestIDFromContext(ctx))
	}
	return errorDetails
}
