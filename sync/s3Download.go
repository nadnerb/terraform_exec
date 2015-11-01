package sync

import (
	"log"
	"os"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//type S3Config struct {
	//Region string
	//S3_bucket string
	//S3_key string
	//Filename string
//}

func Download(config *AwsConfig) {
	sess := session.New(&aws.Config{Region: aws.String(config.Region)})

	client := s3.New(sess)
	result, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(config.S3_bucket),
		Key:    aws.String(config.S3_key),
	})
	if err != nil {
		log.Fatal("Failed to get object", err)
	}

	file, err := os.Create(config.Filename)
	if err != nil {
		log.Fatal("Failed to create file", err)
	}
	if n, err := io.Copy(file, result.Body); err != nil {
		log.Println(n)
		log.Fatal("Failed to copy object to file", err)
	}
	result.Body.Close()
	file.Close()
}

