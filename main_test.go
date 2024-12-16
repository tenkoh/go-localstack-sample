package main

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/docker/go-connections/nat"
	"github.com/google/go-cmp/cmp"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

const (
	region = "ap-northeast-1"
)

func initBucket(
	t *testing.T,
	ctx context.Context,
	client *s3.Client,
	bucket, localPath, s3Key string,
) error {
	t.Helper()

	client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})

	f, _ := os.Open(localPath)
	defer f.Close()

	client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s3Key),
		Body:   f,
	})
	return nil
}

func TestRepository_GetUsers(t *testing.T) {
	// 説明のわかりやすさのためテーブルドリブンテストにはしません
	// 参考：https://golang.testcontainers.org/modules/localstack/#__tabbed_6_2
	t.Setenv("AWS_ACCESS_KEY_ID", "dummy")     // 何かしらの値が必要
	t.Setenv("AWS_SECRET_ACCESS_KEY", "dummy") // 何かしらの値が必要
	ctx := context.Background()

	c, _ := localstack.Run(ctx, "localstack/localstack:3.7.2")
	defer func(c *localstack.LocalStackContainer) {
		testcontainers.TerminateContainer(c)
	}(c)

	provider, _ := testcontainers.NewDockerProvider()
	defer provider.Close()

	// 立ち上がったLocalStackコンテナのエンドポイントを割り出す
	host, _ := provider.DaemonHost(ctx)
	port, _ := c.MappedPort(ctx, nat.Port("4566/tcp"))
	awsEndpoint := "http://" + host + ":" + port.Port()

	awsCfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(awsEndpoint)
	})
	if err := initBucket(t, ctx, client, "test-bucket", "./dev/users.json", "users.json"); err != nil {
		t.Fatal(err)
	}

	repository := NewRepository(client)
	users := repository.GetUsers(ctx, "test-bucket", "users.json")

	want := []User{{ID: 1}}
	if diff := cmp.Diff(want, users); diff != "" {
		t.Error(diff)
	}
}
