package main

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"time"
)

const (
	endpoint        = "minio:9000"
	accessKeyID     = "B9M320QXHD38WUR2MIY3"
	secretAccessKey = "Xv3IrZ0OJt6BA6zqSVac5VpdTZiMQybXCp1YzI20"
	useSSL          = false
	stsDuration     = time.Hour * 1
	bucketName      = "d3invitation"
)

func init() {
	// create bucket
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Succesfully connect to MinIO: %s\n", endpoint)

	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket '%s' already exists.\n", bucketName)
		} else {
			log.Fatalf("Create Bucket '%s' failed: %s\n", bucketName, err)
		}
	} else {
		log.Printf("Create Bucket '%s' successfully。\n", bucketName)
	}

	// flag
	flagBucketName := "flag"
	err = minioClient.MakeBucket(context.Background(), flagBucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), flagBucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket '%s' already exists.\n", flagBucketName)
		} else {
			log.Fatalf("Create Bucket '%s' failed: %s\n", flagBucketName, err)
		}
	} else {
		log.Printf("Create Bucket '%s' successfully.\n", flagBucketName)
	}

	// upload flag
	_, err = minioClient.FPutObject(context.Background(), flagBucketName, "flag", "/flag", minio.PutObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	err = os.Remove("/flag")
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Printf("Removed /flag\n")
	}
}
