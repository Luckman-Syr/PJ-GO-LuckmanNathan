package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type responseError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func Unauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, responseError{
		Error:   "Unauthorized",
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, responseError{
		Error:   "Bad Request",
		Message: message,
	})
}

func NotFound(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusNotFound, responseError{
		Error:   "Not Found",
		Message: message,
	})
}

type responseOK struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Ok(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, responseOK{
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, responseOK{
		Message: message,
		Data:    data,
	})
}

func JustMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
