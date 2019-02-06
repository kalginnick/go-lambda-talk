package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kalginnick/go-lambda-talk/pkg/client"
)

func main() {
	h := handler{
		url:    os.Getenv("DOWNLOAD_URL"),
		bucket: os.Getenv("BUCKET"),
		file:   os.Getenv("FILENAME"),
		region: os.Getenv("AWS_REGION"),
	}

	lambda.Start(h.handle)
}

type handler struct {
	url    string
	bucket string
	file   string
	region string
}

func (h *handler) handle(event events.CloudWatchEvent) error {
	response, err := http.Get(h.url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("response %s from %s", response.Status, h.url)
	}

	return client.WriteS3(h.region, h.bucket, h.file, response.Body)
}
