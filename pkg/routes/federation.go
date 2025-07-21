package routes

import (
	"encoding/json"
	"fmt"
	"hdruk/federated-metadata/pkg"
	"hdruk/federated-metadata/pkg/secrets"
	"hdruk/federated-metadata/pkg/utils"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateFederationHandler Creates an instance of a federation
// secret in gcloud
func CreateFederationHandler(c *gin.Context) {
	method_name := utils.MethodName(0)
	slog.Debug(
		"Creating federation", 
		"x-request-session-id", c.GetHeader("x-request-session-id"),
		"method_name", method_name,
	)
	decoder := json.NewDecoder(c.Request.Body)
	var cs pkg.CreateSecretRequest

	err := decoder.Decode(&cs)
	if err != nil {
		slog.Debug(
			fmt.Sprintf("unable to decode request body: %s", err.Error()),
			"x-request-session-id", c.GetHeader("x-request-session-id"),
			"method_name", method_name,
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unable to decode request body",
			"error":   err.Error(),
		})
		return
	}

	secretCtx := secrets.NewSecrets("", "")
	resp, err := secretCtx.CreateSecret(cs.Path, cs.SecretID, cs.Payload)
	fmt.Printf("%v\n",cs)
	if err != nil {
		slog.Debug(
			fmt.Sprintf("unable to create new secret instance: %s", err.Error()),
			"x-request-session-id", c.GetHeader("x-request-session-id"),
			"method_name", method_name,
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to create new secret instance",
			"test": cs.Path,
			"ID": cs.SecretID,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp,
	})
}

// CreateFederationHandler Creates an instance of a federation
// secret in gcloud
func UpdateFederationHandler(c *gin.Context) {
	method_name := utils.MethodName(0)
	slog.Debug(
		"Updating federation", 
		"x-request-session-id", c.GetHeader("x-request-session-id"),
		"method_name", method_name,
	)

	decoder := json.NewDecoder(c.Request.Body)
	var cs pkg.CreateSecretRequest

	err := decoder.Decode(&cs)
	if err != nil {
		slog.Debug(
			fmt.Sprintf("unable to decode request body: %s", err.Error()),
			"x-request-session-id", c.GetHeader("x-request-session-id"),
			"method_name", method_name,
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unable to decode request body",
			"error":   err.Error(),
		})
		return
	}

	secretCtx := secrets.NewSecrets("", "")
	resp, err := secretCtx.UpdateSecret(cs.Path, cs.SecretID, cs.Payload)
	fmt.Printf("%v\n",cs)
	if err != nil {
		slog.Debug(
			fmt.Sprintf("unable to create new secret instance: %s", err.Error()), 
			"x-request-session-id", c.GetHeader("x-request-session-id"),
			"method_name", method_name,
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to create new secret instance",
			"test": cs.Path,
			"ID": cs.SecretID,
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
	method_name := utils.MethodName(0)
	slog.Debug(
		"Deleting federation", 
		"x-request-session-id", c.GetHeader("x-request-session-id"),
		"method_name", method_name,
	)

	decoder := json.NewDecoder(c.Request.Body)
	var ds pkg.DeleteSecretRequest

	err := decoder.Decode(&ds)
	if err != nil {
		slog.Debug(
			fmt.Sprintf("unable to decode request body: %s", err.Error()),
			"x-request-session-id", c.GetHeader("x-request-session-id"),
			"method_name", method_name,
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unable to decode request body",
			"error":   err.Error(),
		})
	}

	secretCtx := secrets.NewSecrets("", "")
	err = secretCtx.DeleteSecret(ds.SecretID)
	if err != nil {
		slog.Debug(
			fmt.Sprintf("unable to create new secret instance: %s", err.Error()),
			"x-request-session-id", c.GetHeader("x-request-session-id"),
			"method_name", method_name,
		)
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
