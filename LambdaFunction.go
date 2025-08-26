package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event events.SecretsManagerSecretRotationEvent) error {
	fmt.Printf("Rotating secret %s with token %s\n", event.SecretID, event.ClientRequestToken)

	ccRoller, err := NewCCApiKeyRoller(ctx)
	if err != nil {
		fmt.Printf("Failed to Create CC API Key Roller")
		return err
	}

	err = ccRoller.CheckSecret(ctx, &event.SecretID)
	if err != nil {
		fmt.Printf("Secret Validation Failed")
		return err
	}

	//Work through the 4 steps that AWS Secrets Manager requires
	switch event.Step {
	case "createSecret":
		err = ccRoller.CreateSecret(ctx, &event.SecretID, &event.ClientRequestToken)
		if err != nil {
			fmt.Printf("Secret Creation Failed")
			return err
		}
	case "finishSecret":
		err = ccRoller.FinishSecret(ctx, &event.SecretID)
		if err != nil {
			fmt.Printf("Secret Finalization Failed")
			return err
		}
	case "setSecret":
		err = ccRoller.SetSecret(ctx, &event.SecretID)
		if err != nil {
			fmt.Printf("Secret Set Failed")
			return err
		}
	case "testSecret":
		err = ccRoller.TestSecret(ctx, &event.SecretID)
		if err != nil {
			fmt.Printf("Secret Test Failed")
			return err
		}
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
