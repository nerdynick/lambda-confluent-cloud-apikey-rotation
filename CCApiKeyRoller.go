package main

import (
	"context"
	"errors"
	"fmt"

	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"nerdynick.com/lambda-cc-apikey-rotation/services"
)

type CCApiKeyRoller struct {
	awsServ *services.AWSSecretsManagerService
	// ccapi   *ccapi.APIClient
}

func (lb *CCApiKeyRoller) CheckSecret(ctx context.Context, secretID *string) error {
	//Describe the Secret we are working with
	//Both helps verify our permissions and that the secret exists. As well as gets metadata we further neeed to check
	describe := sm.DescribeSecretInput{SecretId: secretID}
	awsSecret, err := lb.awsServ.SM().DescribeSecret(ctx, &describe)
	if err != nil {
		return err
	}

	if awsSecret.RotationEnabled == nil || !*awsSecret.RotationEnabled {
		fmt.Printf("Secret %s is not enabled for rotation\n", *secretID)
		return errors.New("secret is not enabled for rotation")
	}

	return nil
}

/*
From AWS Docs
This method first checks for the existence of a secret for the passed in token. If one does not exist, it will generate a
new secret and put it with the passed in token.
*/
func (lb *CCApiKeyRoller) CreateSecret(ctx context.Context, secretID *string, versionId *string) error {
	_, err := lb.awsServ.SMGetCurrentValue(ctx, secretID, versionId)
	if err != nil {
		return err
	}

	_, err = lb.awsServ.SMGetPendingValue(ctx, secretID, versionId)
	if err != nil {
		//Because the Get Current Value succedded and this failed. We are able to assume that the secret exists but doesn't yet have a pending state

	} else {
		fmt.Printf("Secret %s already has a pending version, skipping creating pending.\n", *secretID)
	}

	return nil
}

func (lb *CCApiKeyRoller) FinishSecret(ctx context.Context, secretID *string) error {
	return nil
}

// Not Used for Standard Implementation as Kafka doesn't support a native way to push credential updates
/*
	From AWS Docs
	This method should set the AWSPENDING secret in the service that the secret belongs to. For example, if the secret is a database
	credential, this method should take the value of the AWSPENDING secret and set the user's password to this value in the database.
*/
func (lb *CCApiKeyRoller) SetSecret(ctx context.Context, secretID *string) error {
	return nil
}

// TODO: Currently requires knowing resource type in order to know what APIs we are testing
/*
	From AWS Docs
	This method should validate that the AWSPENDING secret works in the service that the secret belongs to. For example, if the secret
	is a database credential, this method should validate that the user can login with the password in AWSPENDING and that the user has
	all of the expected permissions against the database.
*/
func (lb *CCApiKeyRoller) TestSecret(ctx context.Context, secretID *string) error {
	return nil
}

func NewCCApiKeyRoller(ctx context.Context) (*CCApiKeyRoller, error) {
	awsServ, err := services.NewAWSSecretsManagerService(ctx)
	if err != nil {
		return nil, err
	}

	ccKR := CCApiKeyRoller{
		awsServ: awsServ,
	}
	return &ccKR, nil
}
