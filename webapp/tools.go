package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
)

type stsCredentialsValue struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

func generateSTSCredentials(objectName string) (*stsCredentialsValue, error) {

	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": ["s3:GetObject", "s3:PutObject"],
				"Resource": ["arn:aws:s3:::%s/%s"]
			}
		]
	}`, bucketName, objectName)

	log.Printf("Policy:\n%s\n", policy)

	stsCredentialProvider, err := credentials.NewSTSAssumeRole("http://minio:9000", credentials.STSAssumeRoleOptions{
		AccessKey:       accessKeyID,
		SecretKey:       secretAccessKey,
		SessionToken:    "",
		Policy:          policy,
		DurationSeconds: int(stsDuration.Seconds()),
	})
	if err != nil {
		return nil, err
	}

	tempCredsValue, err := stsCredentialProvider.Get()
	if err != nil {
		return nil, err
	}

	log.Println("successfully generated STS Credentials!")
	log.Printf("  Access Key ID: %s\n", tempCredsValue.AccessKeyID)
	log.Printf("  Secret Access Key: %s\n", tempCredsValue.SecretAccessKey)
	log.Printf("  Session Token: %s\n", tempCredsValue.SessionToken)

	return &stsCredentialsValue{tempCredsValue.AccessKeyID, tempCredsValue.SecretAccessKey, tempCredsValue.SessionToken}, nil
}

func getObject(stsCredsValue *stsCredentialsValue, bucketName, objectName string) (*minio.Object, error) {
	tmpCreds := credentials.NewStaticV4(stsCredsValue.AccessKeyID, stsCredsValue.SecretAccessKey, stsCredsValue.SessionToken)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  tmpCreds,
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	object, err := minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return object, nil
}

func putObject(stsCredsValue *stsCredentialsValue, bucketName, objectName string, object io.Reader) error {
	tmpCreds := credentials.NewStaticV4(stsCredsValue.AccessKeyID, stsCredsValue.SecretAccessKey, stsCredsValue.SessionToken)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  tmpCreds,
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	_, err = minioClient.PutObject(
		context.Background(),
		bucketName,
		objectName,
		object,
		-1,
		minio.PutObjectOptions{},
	)
	if err != nil {
		return err
	}

	return nil
}
