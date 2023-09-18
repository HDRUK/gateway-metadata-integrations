package secrets

import (
	"context"
	"encoding/json"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// Secrets Defines the shape of a gcloud secrets object
type Secrets struct {
	Parent  string
	Version string
}

// Token Defines the shape of a gcloud secrets object response
type Token struct {
	BearerToken string `json:"bearer_token"`
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
func (s *Secrets) GetSecret() (Token, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return Token{}, fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: s.Version,
	}

	res, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return Token{}, fmt.Errorf("failed to access secret version: %v", err)
	}

	var token Token
	json.Unmarshal(res.Payload.Data, &token)

	return token, nil
}
