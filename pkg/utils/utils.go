package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// HandleError Global error handler utility function
func HandleError(message string, returnVal any) (any, error) {
	return returnVal, fmt.Errorf("%s", message)
}

// FormResponse Helper function to avoid duplication. Forms a gin.H response
// to return
func FormResponse(status int, ret bool, title string, err string) gin.H {
	return gin.H{
		"status":  status,
		"success": ret,
		"title":   title,
		"errors":  err,
	}
}

// IsSuccessfulStatusCode Helper function to determine successful http call
// responses
func IsSuccessfulStatusCode(status int) bool {
	return status >= 200 && status < 300
}
