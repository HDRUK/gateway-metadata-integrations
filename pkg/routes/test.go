package routes

import (
	"encoding/json"
	"fmt"
	"hdruk/federated-metadata/pkg"
	"hdruk/federated-metadata/pkg/pull"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TestFederationHandler will single handidly test each part of operation to
// determine a successful federation configuration.
func TestFederationHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var fed pkg.Federation

	err := decoder.Decode(&fed)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "unable to decode request body",
			"title":   "unabled to decode request body",
			"errors":  err.Error(),
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
		fed.AuthSecretKey,
		fed.AuthType,
		false,
	)

	status, ret, title, err := p.TestCredentials()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"message": ret,
			"title":   title,
			"errors":  err,
		})
		return
	}

	status, ret, title, err = p.TestDatasetsEndpoint()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"message": ret,
			"title":   title,
			"errors":  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  status,
		"message": ret,
		"title":   title,
		"errors":  err,
	})
	return
}
