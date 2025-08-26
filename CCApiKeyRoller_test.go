package main

import (
	"context"
	"testing"
)

func TestNewCCApiKeyRoller(t *testing.T) {
	_, err := NewCCApiKeyRoller(context.Background())
	if err != nil {
		t.Fatalf("Failed to create CC Api Key Roller: %v", err)
	}
}

func TestCheckSecret(t *testing.T) {
	ccKR, err := NewCCApiKeyRoller(context.Background())
	if err != nil {
		t.Fatalf("Failed to create CC Api Key Roller: %v", err)
	}

	var arn = "arn:aws:secretsmanager:us-west-2:648775253477:secret:TestRotationKey-V3lrmM"
	ccKR.CheckSecret(context.Background(), &arn)
}

func TestCreateSecret(t *testing.T) {
	ccKR, err := NewCCApiKeyRoller(context.Background())
	if err != nil {
		t.Fatalf("Failed to create CC Api Key Roller: %v", err)
	}

	var arn = "arn:aws:secretsmanager:us-west-2:648775253477:secret:TestRotationKey-V3lrmM"
	var version = "cc7841b4-fcd1-44cc-acb9-0c74ca0038f7"
	ccKR.CreateSecret(context.Background(), &arn, &version)
}
