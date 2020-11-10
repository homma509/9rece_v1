package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/homma509/9rece/backend/registry"
	"go.uber.org/zap"
)

func handler(ctx context.Context, event events.S3Event) error {
	var logger = registry.Creater().Logger()
	// lc, _ := lambdacontext.FromContext(ctx)
	logger.Info(
		"start lambda function",
		zap.Any("S3Event", event),
		// "CognitoIdentityID", lc.Identity.CognitoIdentityID,
		// "CognitoIdentityPoolID", lc.Identity.CognitoIdentityPoolID,
	)

	err := registry.Creater().UkeController().Move(ctx, event)
	if err != nil {
		logger.Sugar().Errorw(
			"error lambda function",
			zap.String("Result", "error"),
			zap.Error(err),
		)
	}

	logger.Info(
		"end lambda function",
		zap.String("Result", "successful"),
	)

	return err
}

func main() {
	lambda.Start(handler)
}
