package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kalginnick/go-lambda-talk/pkg/client"
)

func main() {
	h := handler{
		url:  os.Getenv("UPLOAD_URL"),
		user: os.Getenv("USER"),
		pswd: os.Getenv("PSWD"),
	}

	lambda.Start(h.handle)
}

type handler struct {
	url  string
	user string
	pswd string
}

func (h *handler) handle(event events.S3Event) error {
	for _, record := range event.Records {
		err := func() error {
			file, err := client.ReadS3(record.AWSRegion, record.S3.Bucket.Name, record.S3.Object.Key)
			if err != nil {
				return err
			}
			defer file.Close()

			return client.WriteFTP(h.user, h.pswd, h.url, record.S3.Object.Key, file)
		}()
		if err != nil {
			return err
		}
	}

	return nil
}
