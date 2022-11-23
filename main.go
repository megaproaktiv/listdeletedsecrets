package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/aws/smithy-go/transport/http"
	"strings"
)

var client *secretsmanager.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	cfg.APIOptions = append(cfg.APIOptions, func(stack *middleware.Stack) error {
		// Attach the custom middleware to the beginning of the Build step
		return stack.Build.Add(secretParameter, middleware.Before)
	})
	client = secretsmanager.NewFromConfig(cfg)

}

var secretParameter = middleware.BuildMiddlewareFunc("IncludeDeleted", func(
	ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler,
) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {

	typedRequest := in.Request.(*http.Request)

	// Add undocumented Parameter
	// var parameter = `{"SortOrder":"desc", "IncludeDeleted": true}`
	var parameter = `
{
  "MaxResults": 100,
  "IncludeDeleted": true,
  "SortOrder": "desc",
  "Filters": []
}
`
	r := strings.NewReader(parameter)
	stream, err := typedRequest.SetStream(r)
	if err != nil {
		panic(err)
	}
	in.Request = stream

	return next.HandleBuild(ctx, in)
})

func main() {
	parms := &secretsmanager.ListSecretsInput{
		SortOrder: types.SortOrderTypeDesc,
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
