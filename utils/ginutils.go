package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/letcommerce/common-module/response"
	"net/http"
	"strconv"
)

// GetIntParam method binds new int Param from ctx and return http.StatusBadRequest if it couldn't parse
func GetIntParam(ctx *gin.Context, paramName string) (int, error) {
	paramVal := ctx.Params.ByName(paramName)
	intParam, err := strconv.Atoi(paramVal)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponseF(err, "can't bind param: %v to int (value = %v)", paramName, paramVal))
	}
	return intParam, err
}

// GetUIntParam method binds new uint Param from ctx and return http.StatusBadRequest if it couldn't parse
func GetUIntParam(ctx *gin.Context, paramName string) (uint, error) {
	paramVal := ctx.Params.ByName(paramName)
	intParam, err := strconv.Atoi(paramVal)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponseF(err, "can't bind param: %v to uint (value = %v)", paramName, paramVal))
	}
	return uint(intParam), err
}

// GetStringParam method binds new string Param from ctx
func GetStringParam(ctx *gin.Context, paramName string) string {
	paramVal := ctx.Params.ByName(paramName)
	return paramVal
}

// BindDTO method binds new DTO from ctx body
func BindDTO[T any](ctx *gin.Context, dto T) (T, error) {
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
func BindAndValidateDTO[T IValidatable](ctx *gin.Context, dto T) (T, error) {
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

func ReturnMessageResponseOrError(ctx *gin.Context, message string, errMessage string, err error) {
	if err == nil {
		ctx.JSON(http.StatusOK, response.Response{Message: message})
	} else {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse(errMessage, err))
	}
}

func ReturnMessageResponseWithIdOrError(ctx *gin.Context, message string, id uint, errMessage string, err error) {
	if err == nil {
		ctx.JSON(http.StatusOK, response.Response{Message: message, ID: id})
	} else {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse(errMessage, err))
	}
}

func CopyDTO[T any](ctx *gin.Context, to T, from interface{}, ignoreEmpty bool) (result T, err error) {
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
