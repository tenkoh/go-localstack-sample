#!/bin/bash
export AWS_ACCESS_KEY_ID=dummy AWS_SECRET_ACCESS_KEY=dummy

awslocal s3 mb s3://test-bucket
awslocal s3 cp /docker-entrypoint-initaws.d/users.json s3://test-bucket/users.json
