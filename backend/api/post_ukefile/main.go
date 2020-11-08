package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/homma509/9rece/backend/logger"
	"github.com/homma509/9rece/backend/registry"
)

func handler(ctx context.Context, event events.S3Event) error {
	// lc, _ := lambdacontext.FromContext(ctx)
	logger.Infof("%v: %#v",
		"S3Event", event,
		// "CognitoIdentityID", lc.Identity.CognitoIdentityID,
		// "CognitoIdentityPoolID", lc.Identity.CognitoIdentityPoolID,
	)

	err := registry.Creater().UkeController().Move(ctx, event)
	if err != nil {
		logger.Errorf("%v: %v\t%v: %v\t%v: %#v",
			"Msg", "couldn't move uke file",
			"Error", err,
		)
	}

	return err
}

func main() {
	lambda.Start(handler)
}
