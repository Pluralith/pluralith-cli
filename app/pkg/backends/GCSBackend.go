package backends

import (
	"fmt"
	"pluralith/pkg/ux"
)

type GCSBackendConfig struct {
	AccessToken                        string `json:"access_token"`
	Bucket                             string `json:"bucket"`
	Credentials                        string `json:"credentials"`
	EncryptionKey                      string `json:"encryption_key"`
	ImpersonateServiceAccount          string `json:"impersonate_service_account"`
	ImpersonateServiceAccountDelegates string `json:"impersonate_service_account_delegates"`
	Prefix                             string `json:"prefix"`
}

func PushToGCSBackend(config TerraformState) error {
	functionName := "PushToGCSBackend"

	uploadSpinner := ux.NewSpinner("Uploading To GCS State Backend", "Diagram Uploaded To GCS State Backend", "Diagram Upload Failed!", true)
	uploadSpinner.Start()

	gcsConfig := GCSBackendConfig{}
	backendErr := MapBackendConfig(config, &gcsConfig)
	if backendErr != nil {
		return fmt.Errorf("loading gcs backend information failed -> %v: %w", functionName, backendErr)
	}

	uploadSpinner.Success()
	return nil
}
