package main

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ockendenjo/beerienteering/pkg/env"
	"github.com/ockendenjo/handler"
)

type H = handler.Handler[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse]

func main() {
	handler.BuildAndStart(func(awsConfig aws.Config) H {
		goLiveStr := handler.MustGetEnv("GO_LIVE_TIME")
		goLiveTime, err := time.Parse(time.RFC3339, goLiveStr)
		if err != nil {
			panic(err)
		}

		h := &lambdaHandler{
			s3Client:   s3.NewFromConfig(awsConfig),
			bucketName: handler.MustGetEnv("BUCKET_NAME"),
			previewKey: handler.MustGetEnv("PREVIEW_KEY"),
			liveKey:    env.OptStr("LIVE_OBJECT_KEY", "2025.json"),
			demoKey:    env.OptStr("DEMO_OBJECT_KEY", "demo.json"),
			goLiveTime: goLiveTime,
		}

		return h.handle
	})
}

type lambdaHandler struct {
	s3Client   *s3.Client
	bucketName string
	previewKey string
	liveKey    string
	demoKey    string
	goLiveTime time.Time
}

func getBody(res *s3.GetObjectOutput, err error) ([]byte, error) {
	if err == nil {
		return io.ReadAll(res.Body)
	}
	if _, ok := errors.AsType[*types.NoSuchKey](err); ok {
		return []byte(`{"demo":false,"stashes":[]}`), nil
	}
	return nil, err
}

func (h *lambdaHandler) handle(ctx *handler.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	res, err := h.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &h.bucketName,
		Key:    new(h.getKey(ctx, event.QueryStringParameters)),
	})

	b, err := getBody(res, err)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(b), Headers: map[string]string{"Content-Type": "application/json"}}, nil
}

func (h *lambdaHandler) getKey(ctx *handler.Context, parameters map[string]string) string {
	logger := ctx.GetLogger()
	if time.Now().After(h.goLiveTime) {
		logger.Info("File is live")
		return h.liveKey
	}
	if parameters["key"] == h.previewKey {
		logger.Info("Preview key supplied")
		return h.liveKey
	}

	logger.Info("Event is not yet live")
	return h.demoKey
}
