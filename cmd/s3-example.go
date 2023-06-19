package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	// Loads AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and AWS_REGION from the environment
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(os.Getenv("AWS_OBJECT_KEY")),
	})
	if err != nil {
		log.Fatalf("unable to download item %v", err)
	}
	defer resp.Body.Close()

	// Read the content
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("unable to read item body %v", err)
	}

	fmt.Printf("Contents:\n%s\nm", string(body))

	// Code to modify
	// Modify the content and create upper case backup
	newContent := strings.ToUpper(
		string(body),
	)

	fmt.Printf("After:\n%s\nm", newContent)
	// Upload the new content back to the S3 bucket
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),      // Bucket name
		Key:    aws.String(os.Getenv("AWS_OBJECT_KEY2")), // Key, including prefix
		Body:   bytes.NewReader([]byte(newContent)),      // the new content to upload
	})
	if err != nil {
		log.Fatalf("unable to upload item %v", err)
	}

	fmt.Println("Successfully uploaded the new content to the S3 bucket")
}
