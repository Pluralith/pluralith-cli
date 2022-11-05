package backends

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureBackendConfig struct {
	AccessKey                 string `json:"access_key"`
	ClientCertificatePassword string `json:"client_certificate_password"`
	ClientCertificatePath     string `json:"client_certificate_path"`
	ClientID                  string `json:"client_id"`
	ClientSecret              string `json:"client_secret"`
	ContainerName             string `json:"container_name"`
	Endpoint                  string `json:"endpoint"`
	Environment               string `json:"environment"`
	Key                       string `json:"key"`
	MetadataHost              string `json:"metadata_host"`
	MsiEndpoint               string `json:"msi_endpoint"`
	ResourceGroupName         string `json:"resource_group_name"`
	SasToken                  string `json:"sas_token"`
	Snapshot                  string `json:"snapshot"`
	StorageAccountName        string `json:"storage_account_name"`
	SubscriptionID            string `json:"subscription_id"`
	TenantID                  string `json:"tenant_id"`
	UseAzureadAuth            string `json:"use_azuread_auth"`
	UseMsi                    string `json:"use_msi"`
}

func PushToAzureBackend(config TerraformState) error {
	functionName := "PushToAzureBackend"

	uploadSpinner := ux.NewSpinner("Uploading To Azure State Backend", "Diagram Uploaded To Azure State Backend", "Diagram Upload Failed!", true)
	uploadSpinner.Start()

	azureConfig := AzureBackendConfig{}
	backendErr := MapBackendConfig(config, &azureConfig)
	if backendErr != nil {
		uploadSpinner.Fail()
		return fmt.Errorf("loading azure backend information failed -> %v: %w", functionName, backendErr)
	}

	url := "https://" + azureConfig.StorageAccountName + ".blob.core.windows.net/" //replace <StorageAccountName> with your Azure storage account name
	ctx := context.Background()

	credential, credentialErr := azidentity.NewDefaultAzureCredential(nil)
	if credentialErr != nil {
		uploadSpinner.Fail()
		return fmt.Errorf("invalid credentials with error -> %v: %w", functionName, credentialErr)
	}

	blobClient, clientErr := azblob.NewClient(url, credential, nil)
	if clientErr != nil {
		uploadSpinner.Fail()
		return fmt.Errorf("creating new azure blob client failed -> %v: %w", functionName, clientErr)
	}

	diagram, diagramErr := os.ReadFile(filepath.Join(auxiliary.StateInstance.WorkingPath, "Infrastructure_Diagram.pdf"))
	if diagramErr != nil {
		uploadSpinner.Fail()
		return fmt.Errorf("reading diagram from disk failed -> %v: %w", functionName, diagramErr)
	}

	keyPath := filepath.Dir(azureConfig.Key)
	diagramPath := filepath.Join(keyPath, "Infrastructure_Diagram.pdf")

	_, uploadErr := blobClient.UploadBuffer(ctx, azureConfig.ContainerName, diagramPath, diagram, nil)
	if uploadErr != nil {
		uploadSpinner.Fail()
		return fmt.Errorf("uploading diagram to azure blob storage failed -> %v: %w", functionName, uploadErr)
	}

	uploadSpinner.Success()
	return nil
}
