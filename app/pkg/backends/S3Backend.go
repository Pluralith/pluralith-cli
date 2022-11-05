package backends

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3BackendConfig struct {
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

func PushToS3Backend(config TerraformState) error {
	functionName := "PushToS3Backend"

	uploadSpinner := ux.NewSpinner("Uploading To S3 State Backend", "Diagram Uploaded To S3 State Backend", "Diagram Upload Failed!", true)
	uploadSpinner.Start()

	s3Config := S3BackendConfig{}
	backendErr := MapBackendConfig(config, &s3Config)
	if backendErr != nil {
		uploadSpinner.Fail()
		return fmt.Errorf("loading s3 backend information failed -> %v: %w", functionName, backendErr)
	}

	s3BucketConfig := &aws.Config{Region: aws.String(s3Config.Region)}
	s3Session, s3Error := session.NewSession(s3BucketConfig)
	if s3Error != nil {
		uploadSpinner.Fail()
		return fmt.Errorf("establishing s3 session failed -> %v: %w", functionName, backendErr)
	}

	uploader := s3manager.NewUploader(s3Session)

	diagram, diagramErr := os.ReadFile(filepath.Join(auxiliary.StateInstance.WorkingPath, "Infrastructure_Diagram.pdf"))
	if diagramErr != nil {
		uploadSpinner.Fail()
		return fmt.Errorf("reading diagram from disk failed -> %v: %w", functionName, diagramErr)
	}

	// Store diagram at same key "directory" as state
	keyPath := filepath.Dir(s3Config.Key)
	diagramPath := filepath.Join(keyPath, "Infrastructure_Diagram.pdf")

	s3Input := &s3manager.UploadInput{
		Bucket:      &s3Config.Bucket,
		Key:         &diagramPath,
		Body:        bytes.NewReader(diagram),
		ContentType: aws.String("application/pdf"),
	}

	_, uploadErr := uploader.UploadWithContext(context.Background(), s3Input)
	if uploadErr != nil {
		uploadSpinner.Fail()
		return fmt.Errorf("uploading diagram to s3 bucket failed -> %v: %w", functionName, uploadErr)
	}

	uploadSpinner.Success()
	return nil
}
