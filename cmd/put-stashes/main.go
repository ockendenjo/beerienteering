package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ockendenjo/beerienteering/pkg/env"
	"github.com/ockendenjo/beerienteering/pkg/stash"
	"github.com/ockendenjo/handler"
)

type H = handler.Handler[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse]

func main() {
	handler.BuildAndStart(func(awsConfig aws.Config) H {
		h := &lambdaHandler{
			s3Client:   s3.NewFromConfig(awsConfig),
			bucketName: handler.MustGetEnv("BUCKET_NAME"),
			liveKey:    env.OptStr("LIVE_OBJECT_KEY", "2025.json"),
		}
		return h.handle
	})
}

type lambdaHandler struct {
	s3Client   *s3.Client
	bucketName string
	liveKey    string
}

func (h *lambdaHandler) handle(ctx *handler.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var sf *stash.StashFile
	err := json.Unmarshal([]byte(request.Body), &sf)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	b, err := json.Marshal(sf)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	_, err = h.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &h.bucketName,
		Key:    &h.liveKey,
		Body:   bytes.NewReader(b),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, nil
}
