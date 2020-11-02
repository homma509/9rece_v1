package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/homma509/9rece/backend/registry"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lc, _ := lambdacontext.FromContext(ctx)

	fmt.Println("Context cognito identity id:", lc.Identity.CognitoIdentityID)
	fmt.Println("Context Cognito pool:", lc.Identity.CognitoIdentityPoolID)
	fmt.Println("Request context cognito id:", event.RequestContext.Identity.CognitoIdentityID)
	fmt.Println("Request context cognito id:", event.RequestContext.Identity.CognitoIdentityPoolID)

	return registry.Creater().URLController().Get(ctx, event)
}

func main() {
	lambda.Start(handler)
}
