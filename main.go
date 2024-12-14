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

type Repository struct {
	client *s3.Client
}

func NewRepository(client *s3.Client) *Repository {
	return &Repository{client}
}

func (r *Repository) GetUsers(ctx context.Context, bucket, key string) []User {
	object, _ := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	defer object.Body.Close()

	var users []User
	json.NewDecoder(object.Body).Decode(&users)
	return users
}

func main() {
	// 説明のためにエラーハンドリングは割愛
	ctx := context.Background()
	s3Config := newS3Config()

	config, _ := config.LoadDefaultConfig(ctx, config.WithRegion(s3Config.region))

	client := s3.NewFromConfig(config, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	repository := NewRepository(client)
	users := repository.GetUsers(ctx, s3Config.bucket, s3Config.key)

	fmt.Printf("%+v", users)
}
