package routes

import (
	"encoding/json"
	"fmt"
	"hdruk/federated-metadata/pkg"
	"hdruk/federated-metadata/pkg/pull"
	"hdruk/federated-metadata/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TestFederationHandler will single handidly test each part of operation to
// determine a successful federation configuration.
func TestFederationHandler(c *gin.Context) {
	var response interface{}

	decoder := json.NewDecoder(c.Request.Body)
	var fed pkg.Federation

	err := decoder.Decode(&fed)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormResponse(http.StatusBadRequest,
			false,
			"unable to decode request body",
			err.Error()))
		return
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
		fed.PID,
		fed.AuthType,
		false,
	)

	response = p.TestCredentials()
	if val, ok := response.(gin.H)["errors"]; ok && val != "" {
		c.JSON(http.StatusOK, response)
		return
	}

	response = p.TestDatasetsEndpoint()

	fmt.Printf("Response received: %+v\n", response)
	if val, ok := response.(gin.H)["errors"]; ok && val != "" {
		fmt.Printf("got error back!!!\n")
		c.JSON(http.StatusOK, response)
		return
	}

	c.JSON(http.StatusOK, utils.FormResponse(http.StatusOK, true, "Test Successful", ""))
}
