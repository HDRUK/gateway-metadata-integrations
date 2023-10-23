package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// Secrets Defines the shape of a gcloud secrets object
type Secrets struct {
	Parent  string
	Version string
}

// Token Defines the shape of a gcloud secrets object response
type BearerTokenResponse struct {
	BearerToken string `json:"bearer_token"`
}

type APIKeyResponse struct {
	APIKey       string `json:"api_key"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// NewSecrets Creates a new Secrets object for interfacing with
// gcloud secret manager
func NewSecrets(parent, version string) *Secrets {
	return &Secrets{
		Parent:  parent,
		Version: version,
	}
}

// GetSecret Returns the current secret version for this secrets
// object version reference
func (s *Secrets) GetSecret(authType string) (any, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("%s/secrets/%s/versions/latest", os.Getenv("GOOGLE_APPLICATION_PROJECT_PATH"), s.Parent),
	}

	res, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %v", err)
	}

	switch strings.ToUpper(authType) {
	case "BEARER":
		var token BearerTokenResponse
		json.Unmarshal(res.Payload.Data, &token)
		return token, nil
	case "API_KEY":
		var token APIKeyResponse
		json.Unmarshal(res.Payload.Data, &token)
		return token, nil
	case "NO_AUTH":
		// Do nothing
	}

	return nil, fmt.Errorf("unable to determine auth type")
}

// CreateSecret Attempts to create a new secret on the given `path`,
// determined by `secretID` within gcloud. Returns the path on success
// or an error otherwise.
func (s *Secrets) CreateSecret(parent, secretID, payload string) (string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		// The most likely causes of the error are:
		//	1. Google application credentials failed
		//	2. Secret already exists
		return "", fmt.Errorf("failed to create secretmanager client %v", err)
	}
	defer client.Close()

	fmt.Printf("%s %s", parent, secretID)

	secretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	result, err := client.CreateSecret(ctx, secretReq)
	if err != nil {
		return "", fmt.Errorf("failed to create secret: %v", err)
	}
	secretName := result.Name

	versionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: result.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte(payload),
		},
	}

	_, err = client.AddSecretVersion(ctx, versionReq)
	if err != nil {
		return "", fmt.Errorf("failed to create new secret version: %v", err)
	}

	return secretName, nil
}

// AddSecretVersion Updates a secret to the new `payload` incrementing
// the gcloud secret version. Returns the secret path on success, error
// otherwise.
func (s *Secrets) AddSecretVersion(path string, payload []byte) (string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: path,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	result, err := client.AddSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to add secret version: %v", err)
	}

	fmt.Printf("added secret version: %s\n", result.Name)
	return result.Name, nil
}

// DeleteSecret Attempts to delete a secret from within gcloud
// secrets manager. Returns nil on success, error otherwise
func (s *Secrets) DeleteSecret(secretID string) error {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	req := &secretmanagerpb.DeleteSecretRequest{
		Name: fmt.Sprintf("%s/secrets/%s", os.Getenv("GOOGLE_APPLICATION_PROJECT_PATH"), secretID),
	}

	if err := client.DeleteSecret(ctx, req); err != nil {
		return fmt.Errorf("failed to delete secret: %v", err)
	}

	return nil
}
