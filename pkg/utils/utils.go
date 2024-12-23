package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

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

// StringInSlice helper function that checks if a string is in an array of strings
// to return true or false
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// FindMissingElements helper functions that loops over a list to find elements
// that are not in another list
// returns a string list of missing elements
func FindMissingElements(list1, list2 []string) []string {
	list2Map := make(map[string]bool)
	for _, elem := range list2 {
		list2Map[elem] = true
	}
	missingElements := []string{}
	for _, elem := range list1 {
		if _, exists := list2Map[elem]; !exists {
			missingElements = append(missingElements, elem)
		}
	}
	return missingElements
}

// WriteGatewayAudit Helper function to write logs to the gateway api audit
// log
func WriteGatewayAudit(message, actionType string, actionName string) {
	enabled, err := strconv.Atoi(os.Getenv("GATEWAY_AUDIT_ENABLED"))
	if err != nil {
		enabled = 0 // couldn't read config, so avoid spamming the API
	}

	if enabled != 1 {
		return
	}

	microseconds := time.Now().UnixNano() / 1000

	payload := []byte(
		fmt.Sprintf(`{
			"user_id": %d,
			"team_id": %d,
			"description": "%s",
			"action_type": "%s",
			"action_service": "%s",
			"action_name": "%s",
			"created_at": %d
		}`, -99, -99, message, actionType, "FMA2", actionName, microseconds))

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
