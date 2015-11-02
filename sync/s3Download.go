package sync

import (
	"log"
	"os"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Download(region string, bucket string, key string, toFile string) error {
	sess := session.New(&aws.Config{Region: aws.String(region)})

	client := s3.New(sess)
	result, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Println("Failed to get object, ", err)
		return err
	}
	file, err := os.Create(toFile)
	if err != nil {
		log.Println("Failed to create file, ", err)
		return err
	}
	if n, err := io.Copy(file, result.Body); err != nil {
		log.Println(n)
		log.Println("Failed to copy object to file, ", err)
		return err
	}
	result.Body.Close()
	file.Close()
	return nil
}

