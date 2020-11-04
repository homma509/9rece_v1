package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/homma509/9rece/backend/registry"
)

func handler(ctx context.Context, event events.S3Event) error {
	return registry.Creater().UkeController().Move(ctx, event)
}

func main() {
	lambda.Start(handler)
}
