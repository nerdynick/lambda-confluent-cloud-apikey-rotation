package services

import (
	"context"
	"errors"
	"fmt"
	"os"

	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	ccapi "github.com/confluentinc/ccloud-sdk-go-v2/apikeys/v2"
	"nerdynick.com/lambda-cc-apikey-rotation/models"
)

// ENV Vars for the different ways to bootstrap the CC API Client
const CCClientAWSSecret = "AWS_SECRET"                 //Via AWS Secrets
const CCClientAPIKey = "CONFLUENT_CLOUD_API_KEY"       //Via ENV Var directly
const CCClientAPISecret = "CONFLUENT_CLOUD_API_SECRET" //Via ENV Var directly
const CCUserAgentAdd = "/nn/aws/apikey/lambda-rotation"

type CCService struct {
	cc    *ccapi.APIClient
	ccCtx *context.Context
}

func (ccS *CCService) NewCloudAPIKey(displayName string, desc string, owner string, id string) {
	resource := ccapi.NewObjectReferenceWithDefaults()
	resource.SetId("cloud")

	spec := ccapi.NewIamV2ApiKeySpec()
	spec.SetDescription(desc)
	spec.SetDisplayName(displayName)
	spec.SetResource(*ccapi.NewObjectReference())
	ccS.cc.APIKeysIamV2Api.CreateIamV2ApiKey(*ccS.ccCtx).IamV2ApiKey(*ccapi.NewIamV2ApiKey())
}

func NewCCService(ctx context.Context, awsSM *sm.Client) (*CCService, error) {
	var newCtx context.Context

	// Confluent Creds need to be provided via one of the below supported routes
	awsSecret, awsSecDefined := os.LookupEnv(CCClientAWSSecret) //Used for fetching Creds from AWS Secrets Manager
	apiKey, apiKeyDef := os.LookupEnv(CCClientAPIKey)           //Used to get Creds from ENV
	apiSec, apiSecDef := os.LookupEnv(CCClientAPISecret)        //Used to get Creds from ENV

	if awsSecDefined {
		sec, err := awsSM.GetSecretValue(ctx, &sm.GetSecretValueInput{SecretId: &awsSecret})
		if err != nil {
			fmt.Printf("Failed to get AWS Secret %s to use for Confluent Cloud API", awsSecret)
			return nil, errors.New("Failed to get CC Creds from AWS Secrets Manager")
		}

		creds, err := models.CredentialsFromJSON(sec.SecretBinary)
		if err != nil {
			fmt.Printf("AWS Secret %s is corrupt. Failed to deserilize from JSON", awsSecret)
			return nil, errors.New("AWS Secret is corrupt. Failed to deserilize from JSON")
		}

		newCtx = context.WithValue(ctx, ccapi.ContextAPIKeys, ccapi.APIKey{
			Key: creds.APIKey,
		})
	} else if apiKeyDef && apiSecDef {
		newCtx = context.WithValue(ctx, ccapi.ContextBasicAuth, ccapi.BasicAuth{
			UserName: apiKey,
			Password: apiSec,
		})
	} else {
		fmt.Printf("No Confluent Cloud credentials found")
		return nil, fmt.Errorf("no AWS Secret ('%s') or CC API Key ('%s') and/or CC API Secret ('%s') provided for CC API Credentials", CCClientAWSSecret, CCClientAPIKey, CCClientAPISecret)
	}

	//Create Default Config & Setup UserAgent to be Unique
	var ccCfg = ccapi.NewConfiguration()
	ccCfg.UserAgent = ccCfg.UserAgent + CCUserAgentAdd

	//TODO: Add support for changing out CC API URL

	return &CCService{
		cc:    ccapi.NewAPIClient(ccCfg),
		ccCtx: &newCtx,
	}, nil
}
