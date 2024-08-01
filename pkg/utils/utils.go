package utils

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
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
func WriteGatewayAudit(message, actionType string) {
	enabled, err := strconv.Atoi(os.Getenv("AUDIT_LOG_ENABLED"))
	if err != nil {
		enabled = 0 // couldn't read config, so avoid spamming the API
	}

	if enabled != 1 {
		return
	}

	ctx := context.Background()
	projectId := os.Getenv("PUBSUB_PROJECT_ID")
	topicName := os.Getenv("PUBSUB_TOPIC_NAME")

	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		slog.Info(fmt.Sprintf("Failed to create client: %s", err.Error()))
	}
	defer client.Close()

	payload := []byte(
		fmt.Sprintf(`{
			"user_id": %d,
			"team_id": %d,
			"description": "%s",
			"action_type": "%s",
			"action_service": "%s",
			"created_at": %d
		}`, -99, -99, message, actionType, "FMA2", time.Now().UnixMicro()))

	pubSubMessage := &pubsub.Message{Data: payload}

	topic := client.Topic(topicName)
	res := topic.Publish(ctx, pubSubMessage)

	id, err := res.Get(ctx)
	if err != nil {
		slog.Info(fmt.Sprintln(err.Error()))
	}
	slog.Debug(fmt.Sprintf("Message published, id: %s", id))
}
