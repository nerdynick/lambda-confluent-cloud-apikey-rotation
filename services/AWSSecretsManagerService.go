package services

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const AWSVersionIDCurrent string = "AWSCURRENT"
const AWSVersionIDPending string = "AWSPENDING"
const AWSVersionIDPrevious = "AWSPREVIOUS"

type AWSSecretsManagerService struct {
	awsSM *sm.Client

	//Short Cuts to get around AWS's SDK requiring pointers for everything and GoLang not supporting CONST+Pointers
	versionID_Current  string
	versionID_Pending  string
	versionID_Previous string
}

func (sms *AWSSecretsManagerService) SM() *sm.Client {
	return sms.awsSM
}

func (sms *AWSSecretsManagerService) SMGetCurrentValue(ctx context.Context, secretID *string, versionID *string) (*sm.GetSecretValueOutput, error) {
	return sms.awsSM.GetSecretValue(ctx, &sm.GetSecretValueInput{SecretId: secretID, VersionId: versionID, VersionStage: &sms.versionID_Current})
}

func (sms *AWSSecretsManagerService) SMGetPendingValue(ctx context.Context, secretID *string, versionID *string) (*sm.GetSecretValueOutput, error) {
	return sms.awsSM.GetSecretValue(ctx, &sm.GetSecretValueInput{SecretId: secretID, VersionId: versionID, VersionStage: &sms.versionID_Pending})
}

func NewAWSSecretsManagerService(ctx context.Context) (*AWSSecretsManagerService, error) {
	//Load Up AWS Configs
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	asms := AWSSecretsManagerService{
		awsSM:              sm.NewFromConfig(awsCfg),
		versionID_Current:  AWSVersionIDCurrent,
		versionID_Pending:  AWSVersionIDPending,
		versionID_Previous: AWSVersionIDPrevious,
	}

	return &asms, nil
}
