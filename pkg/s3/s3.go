package s3

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadToR2(bucketName, key, filePath string) (string, error) {
	// Initialize a session that the SDK will use to load credentials from the shared credentials file ~/.aws/credentials
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Endpoint:    aws.String(os.Getenv("R2_ENDPOINT")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	// Open the file for use
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q, %v", filePath, err)
	}
	defer file.Close()

	// Get the file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buffer),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}

	log.Printf("Successfully uploaded %s to %s/%s\n", filePath, bucketName, key)
	return result.Location, nil
}
