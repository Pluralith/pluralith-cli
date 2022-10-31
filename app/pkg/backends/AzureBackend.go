package backends

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
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

func LoadAzureBackendConfig(tfState TerraformState) interface{} {
	azureConfig := AzureBackendConfig{}

	config := &mapstructure.DecoderConfig{TagName: "json"}
	config.Result = &azureConfig

	decoder, _ := mapstructure.NewDecoder(config)
	decoder.Decode(tfState.Backend.Config)

	return azureConfig
}

func InitAzureBackend(tfState TerraformState) {
	azureConfig := LoadAzureBackendConfig(tfState)
	fmt.Println(azureConfig)
	// credential, err := azidentity.NewDefaultAzureCredential(nil)
	// if err != nil {
	// 	log.Fatal("Invalid credentials with error: " + err.Error())
	// }

	// serviceClient, err := azblob.NewClient(url, credential, nil)
	// if err != nil {
	// 	log.Fatal("Invalid credentials with error: " + err.Error())
	// }
}