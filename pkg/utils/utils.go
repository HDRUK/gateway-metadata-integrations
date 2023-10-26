package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

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

// WriteGatewayAudit Helper function to write logs to the gateway api audit
// log
func WriteGatewayAudit(message, actionType string) {
	payload := []byte(
		fmt.Sprintf(`{
			"user_id": %d,
			"team_id": %d,
			"description": "%s",
			"action_type": "%s",
			"action_service": "%s"
		}`, -99, -99, message, actionType, "FMA2"))

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/%s", os.Getenv("GATEWAY_API_URL"), "audit"),
		bytes.NewBuffer(payload),
	)
	if err != nil {
		fmt.Printf("unable to form new request for api audit log entry %v", err.Error())
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("unable to call gateway api for audit log entry %v", err.Error())
	}
	defer res.Body.Close()

	fmt.Printf("%v", payload)
}
