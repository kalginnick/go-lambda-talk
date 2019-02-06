package client

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func ReadS3(region, bucket, file string) (io.ReadCloser, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, err
	}

	loader := s3.New(sess)
	response, err := loader.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file),
	})
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func WriteS3(region, bucket, file string, data io.Reader) error {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return err
	}

	loader := s3manager.NewUploader(sess)
	_, err = loader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file),
		Body:   data,
	})
	if err != nil {
		return err
	}

	return nil
}
