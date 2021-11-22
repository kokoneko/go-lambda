package main

import (
	"context"
	"fmt"
	"app/event"
	
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, events events.S3Event) (string, error) {
	fmt.Println(event.S3lambda(events).Key)
	return event.S3lambda(events).Key, nil
}

func main() {
	lambda.Start(HandleRequest)
}