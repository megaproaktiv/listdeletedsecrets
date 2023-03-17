package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

var client *secretsmanager.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client = secretsmanager.NewFromConfig(cfg)
}

func main() {
	parms := &secretsmanager.ListSecretsInput{
		SortOrder:              types.SortOrderTypeDesc,
		IncludePlannedDeletion: aws.Bool(true),
		MaxResults:             aws.Int32(100),
	}
	resp, err := client.ListSecrets(context.TODO(), parms)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Results\n=======\n")
	for _, s := range resp.SecretList {
		fmt.Printf("Secret: %v / deleted on %v\n", *s.Name, s.DeletedDate)
	}
}
