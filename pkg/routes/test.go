package routes

import (
	"encoding/json"
	"fmt"
	"hdruk/federated-metadata/pkg"
	"hdruk/federated-metadata/pkg/pull"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TestCredentialsHandler Handles the test instruction to validate a
// newly created federation integration
func TestCredentialsHandler(c *gin.Context) {

	decoder := json.NewDecoder(c.Request.Body)
	var fed pkg.Federation

	err := decoder.Decode(&fed)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unable to decode request body",
			"error":   err.Error(),
		})
	}

	// Create a new Pull Object to test this integration
	p := pull.NewPull(
		fed.ID,
		fmt.Sprintf("%s%s", fed.EndpointBaseURL,
			fed.EndpointDatasets),
		fmt.Sprintf("%s%s", fed.EndpointBaseURL,
			fed.EndpointDataset),
		"",
		"",
		// "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiIsImtpZCI6ImYwODZlYjRkLTBkZDMtNGFkZi05Zjk0LTQ1YjhhNjE2ZjRiYSJ9.eyJleHAiOjE3MDIzMzkxOTksIm5iZiI6MTY5NDUyODY5NCwiaWF0IjoxNjk0NTI4Njk0LCJzdWIiOiJhYzkyM2I0ZC02NjE3LTQ5YzEtOTA0ZS05YTg1NDJmY2NlYWIiLCJpc3MiOiJodHRwczovL2ZhaXIucHJldmlldy5hcmlkaGlhLmlvL2FwaSIsImF1ZCI6InByZXZpZXciLCJ0b2tlbiI6ImhkcnVrLXRlc3QtaW50ZWdyYXRpb25zIn0.DLrsiL3iizV743kYsvQLECeh7-loGTE8C2IlOWRmPK8",
		fed.AuthSecretKey,
		fed.AuthType,
		false,
	)

	status, ret, err := p.TestCredentials()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": fmt.Sprintf("%d, %v, %v", status, ret, err),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"response": fmt.Sprintf("%d, %v, %v", status, ret, err),
	})
}

// TestDatasetsEndpointHandler Handles the test instruction to validate a
// newly created federation integration
func TestDatasetsEndpointHandler(c *gin.Context) {

	decoder := json.NewDecoder(c.Request.Body)
	var fed pkg.Federation

	err := decoder.Decode(&fed)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unable to decode request body",
			"error":   err.Error(),
		})
	}

	// Create a new Pull Object to test this integration
	p := pull.NewPull(
		fed.ID,
		fmt.Sprintf("%s%s", fed.EndpointBaseURL,
			fed.EndpointDatasets),
		fmt.Sprintf("%s%s", fed.EndpointBaseURL,
			fed.EndpointDataset),
		"",
		"",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiIsImtpZCI6ImYwODZlYjRkLTBkZDMtNGFkZi05Zjk0LTQ1YjhhNjE2ZjRiYSJ9.eyJleHAiOjE3MDIzMzkxOTksIm5iZiI6MTY5NDUyODY5NCwiaWF0IjoxNjk0NTI4Njk0LCJzdWIiOiJhYzkyM2I0ZC02NjE3LTQ5YzEtOTA0ZS05YTg1NDJmY2NlYWIiLCJpc3MiOiJodHRwczovL2ZhaXIucHJldmlldy5hcmlkaGlhLmlvL2FwaSIsImF1ZCI6InByZXZpZXciLCJ0b2tlbiI6ImhkcnVrLXRlc3QtaW50ZWdyYXRpb25zIn0.DLrsiL3iizV743kYsvQLECeh7-loGTE8C2IlOWRmPK8",
		// fed.AuthSecretKey,
		fed.AuthType,
		false,
	)

	status, ret, err := p.TestDatasetsEndpoint()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": fmt.Sprintf("%d, %v, %v", status, ret, err),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"response": fmt.Sprintf("%d, %v, %v", status, ret, err),
	})
}
