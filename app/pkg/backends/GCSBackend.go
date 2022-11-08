package backends

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"

	"cloud.google.com/go/storage"
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
		uploadSpinner.Fail()
		return fmt.Errorf("loading gcs backend information failed -> %v: %w", functionName, backendErr)
	}

	ctx := context.Background()
	gcsClient, clientErr := storage.NewClient(ctx)
	if clientErr != nil {
		uploadSpinner.Fail()
		return fmt.Errorf("creating gcs client failed -> %v: %w", functionName, clientErr)
	}

	diagram, diagramErr := os.ReadFile(filepath.Join(auxiliary.StateInstance.WorkingPath, "Infrastructure_Diagram.pdf"))
	if diagramErr != nil {
		uploadSpinner.Fail()
		return fmt.Errorf("reading diagram from disk failed -> %v: %w", functionName, diagramErr)
	}

	// Store diagram at same key "directory" as state
	diagramPath := filepath.Join(gcsConfig.Prefix, "Infrastructure_Diagram.pdf")
	gcsObject := gcsClient.Bucket(gcsConfig.Bucket).Object(diagramPath)

	gcsWriter := gcsObject.NewWriter(ctx)
	if _, uploadErr := io.Copy(gcsWriter, bytes.NewReader(diagram)); uploadErr != nil {
		uploadSpinner.Fail()
		return fmt.Errorf("uploading diagram to s3 bucket failed -> %v: %w", functionName, uploadErr)
	}
	if closeErr := gcsWriter.Close(); closeErr != nil {
		uploadSpinner.Fail()
		return fmt.Errorf("uploading diagram to s3 bucket failed -> %v: %w", functionName, closeErr)
	}

	uploadSpinner.Success()
	return nil
}
