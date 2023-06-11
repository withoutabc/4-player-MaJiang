package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type NormSuccess struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
	Data   any    `json:"data"`
}

func RespNormSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, NormSuccess{
		Status: 200,
		Info:   "success",
		Data:   data,
	})
}

type RespTemplate struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
}

func RespOK(c *gin.Context) {
	c.JSON(http.StatusOK, RespTemplate{
		Status: 200,
		Info:   "success",
	})
}

var Unauthorized = RespTemplate{
	Status: 401,
	Info:   "unauthorized",
}

func RespUnauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Unauthorized)
}

var InternalErr = RespTemplate{
	Status: 500,
	Info:   "internal error",
}

func RespInternalErr(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, InternalErr)
}

func RespNormErr(c *gin.Context, errCode int) {
	c.JSON(http.StatusBadRequest, RespTemplate{
		Status: errCode,
		Info:   ErrorCodeMap[errCode].Error(),
	})
}
