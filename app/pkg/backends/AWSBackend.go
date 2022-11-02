package backends

import "fmt"

type AWSBackendConfig struct {
	AccessKey                   string `json:"access_key"`
	ACL                         string `json:"acl"`
	AssumeRoleDurationSeconds   string `json:"assume_role_duration_seconds"`
	AssumeRolePolicy            string `json:"assume_role_policy"`
	AssumeRolePolicyArns        string `json:"assume_role_policy_arns"`
	AssumeRoleTags              string `json:"assume_role_tags"`
	AssumeRoleTransitiveTagKeys string `json:"assume_role_transitive_tag_keys"`
	Bucket                      string `json:"bucket"`
	DynamodbEndpoint            string `json:"dynamodb_endpoint"`
	DynamodbTable               string `json:"dynamodb_table"`
	Encrypt                     string `json:"encrypt"`
	Endpoint                    string `json:"endpoint"`
	ExternalID                  string `json:"external_id"`
	ForcePathStyle              string `json:"force_path_style"`
	IamEndpoint                 string `json:"iam_endpoint"`
	Key                         string `json:"key"`
	KmsKeyID                    string `json:"kms_key_id"`
	MaxRetries                  string `json:"max_retries"`
	Profile                     string `json:"profile"`
	Region                      string `json:"region"`
	RoleArn                     string `json:"role_arn"`
	SecretKey                   string `json:"secret_key"`
	SessionName                 string `json:"session_name"`
	SharedCredentialsFile       string `json:"shared_credentials_file"`
	SkipCredentialsValidation   string `json:"skip_credentials_validation"`
	SkipMetadataAPICheck        string `json:"skip_metadata_api_check"`
	SkipRegionValidation        string `json:"skip_region_validation"`
	SseCustomerKey              string `json:"sse_customer_key"`
	StsEndpoint                 string `json:"sts_endpoint"`
	Token                       string `json:"token"`
	WorkspaceKeyPrefix          string `json:"workspace_key_prefix"`
}

func PushToAWSBackend(config TerraformState) error {
	functionName := "PushToAWSBackend"

	awsConfig := AWSBackendConfig{}
	backendErr := MapBackendConfig(config, &awsConfig)
	if backendErr != nil {
		return fmt.Errorf("loading aws backend information failed -> %v: %w", functionName, backendErr)
	}

	return nil
}
