package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	presignedURL := presignurl()

	// Open the file that you want to upload
	file, err := os.Open("image.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, _ := io.ReadAll(file)

	// Create an HTTP client
	client := http.Client{
		Timeout: 0, // No timeout for the client
	}

	// Create a PUT request to the pre-signed URL
	req, err := http.NewRequest("PUT", presignedURL, strings.NewReader(string(data)))
	if err != nil {
		log.Fatal(err)
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to upload file. Status code: %d", resp.StatusCode)
	}

	log.Println("File uploaded successfully!")
}

func presignurl() string {
	// Create S3 service client
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String("bucket-name"),
		Key:    aws.String("image.jpeg"),
		// Body:   strings.NewReader(""),
	})
	str, err := req.Presign(15 * time.Minute)

	log.Println(str, " err:", err)

	return str
}
