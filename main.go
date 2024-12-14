package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Config struct {
	region string
	bucket string
	key    string
}

func newS3Config() *s3Config {
	return &s3Config{
		region: os.Getenv("AWS_REGION"),
		bucket: os.Getenv("AWS_S3_BUCKET"),
		key:    os.Getenv("AWS_S3_KEY"),
	}
}

type User struct {
	ID int `json:"id"`
}

func main() {
	ctx := context.Background()
	s3Config := newS3Config()

	config, err := config.LoadDefaultConfig(ctx, config.WithRegion(s3Config.region))
	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(config, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	object, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s3Config.bucket),
		Key:    aws.String(s3Config.key),
	})

	if err != nil {
		panic(err)
	}
	defer object.Body.Close()

	var users []User
	if err := json.NewDecoder(object.Body).Decode(&users); err != nil {
		panic(err)
	}

	fmt.Printf("%+v", users)
}
