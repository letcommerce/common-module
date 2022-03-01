package ginutils

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/letcommerce/common-module/response"
	"net/http"
	"strconv"
)

var (
	ctx *gin.Context
)

func Init(ginCtx *gin.Context) {
	ctx = ginCtx
}

// GetIntParam method binds new int Param from ctx and return http.StatusBadRequest if it couldn't parse
func GetIntParam(paramName string) (int, error) {
	paramVal := ctx.Params.ByName(paramName)
	intParam, err := strconv.Atoi(paramVal)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponseF(err, "can't bind param: %v to int (value = %v)", paramName, paramVal))
	}
	return intParam, err
}

// GetUIntParam method binds new uint Param from ctx and return http.StatusBadRequest if it couldn't parse
func GetUIntParam(paramName string) (uint, error) {
	paramVal := ctx.Params.ByName(paramName)
	intParam, err := strconv.Atoi(paramVal)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponseF(err, "can't bind param: %v to uint (value = %v)", paramName, paramVal))
	}
	return uint(intParam), err
}

// GetStringParam method binds new string Param from ctx
func GetStringParam(paramName string) string {
	paramVal := ctx.Params.ByName(paramName)
	return paramVal
}

// BindDTO method binds new DTO from ctx body
func BindDTO[T any](dto T) (T, error) {
	err := ctx.Bind(&dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Got error while binding dto", err))
	}
	return dto, err
}

type IValidatable interface {
	Validate() error
}

// BindAndValidateDTO method binds new DTO from ctx body
func BindAndValidateDTO[T IValidatable](dto T) (T, error) {
	var null T
	err := ctx.Bind(&dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Got error while binding dto", err))
		return null, err
	}
	err = dto.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Got error while validating dto", err))
		return null, err
	}
	return dto, err
}

func ReturnResultOrError(result interface{}, errMessage string, err error) {
	if err == nil {
		ctx.JSON(http.StatusOK, result)
	} else {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse(errMessage, err))
	}
}

func ReturnMessageResponseOrError(message string, errMessage string, err error) {
	if err == nil {
		ctx.JSON(http.StatusOK, response.Response{Message: message})
	} else {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse(errMessage, err))
	}
}

func ReturnMessageResponseWithIdOrError(message string, id uint, errMessage string, err error) {
	if err == nil {
		ctx.JSON(http.StatusOK, response.Response{Message: message, ID: id})
	} else {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse(errMessage, err))
	}
}

func CopyDTO[T any](to T, from interface{}, ignoreEmpty bool) (result T, err error) {
	var null T
	if ignoreEmpty {
		err = copier.CopyWithOption(&to, from, copier.Option{IgnoreEmpty: true})
	} else {
		err = copier.Copy(&to, from)
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Got error while coping", err))
		return null, err
	}
	return to, nil
}
