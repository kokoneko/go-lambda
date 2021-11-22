package main

import (
	"context"

	event "app/event"
	resize "app/resize"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, events events.S3Event) (string, error) {
	s3TrigerInfo := event.S3lambda(events)
	resize.ExecResize(s3TrigerInfo.Bucket, s3TrigerInfo.Key)
	return "", nil
}

func main() {
	lambda.Start(HandleRequest)
}