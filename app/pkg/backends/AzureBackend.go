package backends

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"

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

	azureConfig := AzureBackendConfig{}
	backendErr := MapBackendConfig(config, &azureConfig)
	if backendErr != nil {
		return fmt.Errorf("loading azure backend information failed -> %v: %w", functionName, backendErr)
	}

	fmt.Println(azureConfig)

	fmt.Printf("Azure Blob storage quick start sample\n")

	url := "https://" + azureConfig.StorageAccountName + ".blob.core.windows.net/" //replace <StorageAccountName> with your Azure storage account name
	ctx := context.Background()

	fmt.Println(url)

	// Create a default request pipeline using your storage account name and account key.
	credential, credentialErr := azidentity.NewDefaultAzureCredential(nil)
	if credentialErr != nil {
		return fmt.Errorf("invalid credentials with error -> %v: %w", functionName, credentialErr)
	}

	blobClient, clientErr := azblob.NewClient(url, credential, nil)
	if clientErr != nil {
		return fmt.Errorf("creating new azure blob client failed -> %v: %w", functionName, clientErr)
	}

	diagram, diagramErr := os.ReadFile(filepath.Join(auxiliary.StateInstance.WorkingPath, "Infrastructure_Diagram.pdf"))
	if diagramErr != nil {
		return fmt.Errorf("reading diagram from disk failed -> %v: %w", functionName, diagramErr)
	}

	// Upload to data to blob storage
	_, uploadErr := blobClient.UploadBuffer(ctx, azureConfig.ContainerName, "Infrastructure_Diagram.pdf", diagram, nil)
	if uploadErr != nil {
		return fmt.Errorf("uploading diagram to azure blob storage failed -> %v: %w", functionName, uploadErr)
	}

	return nil
}
