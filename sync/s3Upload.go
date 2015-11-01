package sync

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Upload(config *AwsConfig) {

	session := session.New(&aws.Config{Region: aws.String(config.Region)})
	client := s3.New(session)
	//result, err := client.CreateBucket(&s3.CreateBucketInput{
		//Bucket: &bucket,
	//})
	//if err != nil {
		//log.Println("Failed to create bucket", err)
		//return
	//}

	_, err := client.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader("Hello!"),
		Bucket: &config.S3_bucket,
		Key:    &config.S3_key,
	})
	if err != nil {
		log.Printf("Failed to upload data to %s/%s, %s\n", config.S3_bucket, config.S3_key, err)
		return
	}

	log.Printf("Successfully created bucket %s and uploaded data with key %s\n", config.S3_bucket, config.S3_key);
}

