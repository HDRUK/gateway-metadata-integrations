package routes

import (
	"encoding/json"
	"hdruk/federated-metadata/pkg"
	"hdruk/federated-metadata/pkg/secrets"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateFederationHandler Creates an instance of a federation
// secret in gcloud
func CreateFederationHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var cs pkg.CreateSecretRequest

	err := decoder.Decode(&cs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unable to decode request body",
			"error":   err.Error(),
		})
		return
	}

	secretCtx := secrets.NewSecrets("", "")
	resp, err := secretCtx.CreateSecret(cs.Path, cs.SecretID, cs.Payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to create new secret instance",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp,
	})
}

// TODO Update version handler - not needed right now
// TODO List version handler - not needed right now

// DeleteFederationHandler Attempts to call the gcloud delete secrets
// function to remove a secret held within the system.
func DeleteFederationHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var ds pkg.DeleteSecretRequest

	err := decoder.Decode(&ds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unable to decode request body",
			"error":   err.Error(),
		})
	}

	secretCtx := secrets.NewSecrets("", "")
	err = secretCtx.DeleteSecret(ds.SecretID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to delete secret instance",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
