package secrets

import (
	"context"
	"encoding/json"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/iterator"
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

// CreateSecret Attempts to create a new secret on the given `path`,
// determined by `secretID` within gcloud. Returns the path on success
// or an error otherwise.
func (s *Secrets) CreateSecret(parent, secretID string) (string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		// The most likely causes of the error are:
		//	1. Google application credentials failed
		//	2. Secret already exists
		return "", fmt.Errorf("failed to create secretmanager client %v", err)
	}
	defer client.Close()

	req := &secretmanagerpb.CreateSecretRequest{
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

	result, err := client.CreateSecret(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create secret: %v", err)
	}
	fmt.Printf("created secret: %s\n", result.Name)
	return result.Name, nil
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

// ListSecrets Returns a list of secrets held within gcloud. Does not return
// the encryped payload, only the reference paths for said secrets. Returns
// errors in parallel if any are encountered.
func (s *Secrets) ListSecrets(parent string) (secrets []string, errors []error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return secrets, append(errors, err)
	}
	defer client.Close()

	req := &secretmanagerpb.ListSecretsRequest{
		Parent: parent,
	}

	it := client.ListSecrets(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			errors = append(errors, err)
			secrets = append(secrets, "")
			continue
		}

		secrets = append(secrets, resp.Name)
		errors = append(errors, nil)
	}

	return secrets, errors
}

// DeleteSecret Attempts to delete a secret from within gcloud
// secrets manager. Returns nil on success, error otherwise
func (s *Secrets) DeleteSecret(name string) error {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	req := &secretmanagerpb.DeleteSecretRequest{
		Name: name,
	}

	if err := client.DeleteSecret(ctx, req); err != nil {
		return fmt.Errorf("failed to delete secret: %v", err)
	}

	return nil
}
