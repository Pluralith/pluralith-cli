package backends

import "fmt"

type GCPBackendConfig struct {
	AccessToken                        string `json:"access_token"`
	Bucket                             string `json:"bucket"`
	Credentials                        string `json:"credentials"`
	EncryptionKey                      string `json:"encryption_key"`
	ImpersonateServiceAccount          string `json:"impersonate_service_account"`
	ImpersonateServiceAccountDelegates string `json:"impersonate_service_account_delegates"`
	Prefix                             string `json:"prefix"`
}

func PushToGCPBackend(config TerraformState) error {
	functionName := "PushToGCPBackend"

	gcpConfig := GCPBackendConfig{}
	backendErr := MapBackendConfig(config, &gcpConfig)
	if backendErr != nil {
		return fmt.Errorf("loading gcp backend information failed -> %v: %w", functionName, backendErr)
	}

	return nil
}
