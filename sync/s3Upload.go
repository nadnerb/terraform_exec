package sync

import (
	"bytes"
	"net/http"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Upload(region string, bucket string, key string, filename string) error {

	session := session.New(&aws.Config{Region: aws.String(region)})
	client := s3.New(session)

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal("error opening file\n", err)
		os.Exit(1)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()

	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer) // convert to io.ReadSeeker type
	fileType := http.DetectContentType(buffer)

	_, err = client.PutObject(&s3.PutObjectInput{
		Body:   fileBytes,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	})
	if err != nil {
		log.Printf("Failed to upload data to %s/%s, %s\n", bucket, key, err)
		return err
	}

	//log.Println(awsutil.StringValue(result))
	return nil
}

