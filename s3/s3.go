package s3

import (
	"fmt"
	"log"
	"path/filepath"
	"os"
	"io"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/credentials"
	// "github.com/gorilla/mux"
	"net/http"
)

func createAWSSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Endpoint: aws.String("https://s3.amazonaws.com"),
		Credentials: credentials.NewStaticCredentials("AKIAT3QWW3GJQBHP5JUQ", "XtJJjStvSkVGupofkGXDfikSSub3tGdn1Cjtw3Iu", ""),
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create AWS session: %v", err)
	}
	return sess, nil
}

func CheckIfS3FileExists(key string) (bool, error){

		sess, err := createAWSSession()

		if err !=nil{
			return false, err
		}

		// Create an S3 client using the session
		s3Client := s3.New(sess)

		bucketName := "bhavcopy"
		folderPrefix := key
	
		// Create an input to list objects in the S3 bucket with the specified prefix
		listObjectsInput := &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
			Prefix: aws.String(folderPrefix),
		}
	
		result, err := s3Client.ListObjectsV2(listObjectsInput)
		if err != nil {
			return false, fmt.Errorf("Failed to list objects in S3 bucket: %v", err)
		}
	
		// Check if any objects were found with the specified prefix
		if len(result.Contents) > 0 {
			fmt.Printf("Folder '%s' exists in S3 bucket '%s'\n", folderPrefix, bucketName)
			return true, nil
		} else {
			fmt.Printf("Folder '%s' does not exist in S3 bucket '%s'\n", folderPrefix, bucketName)
			return false, nil
		}
}

func UploadFileToS3(bucketName, folderName, filePath string) error {

	sess, err := createAWSSession()
	if err != nil {
		return err
	}
	svc := s3.New(sess)

	objectKey := fmt.Sprintf("%s/%s", folderName, filepath.Base(filePath))

	// Open the file to upload
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Specify the parameters for the S3 upload
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	}

	// Upload the file to S3
	_, err = svc.PutObject(params)
	if err != nil {
		return err
	}

	return nil
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request, bucketName, key string) error {

	sess, err := createAWSSession()
	if err != nil {
		return err
	}
	s3Client := s3.New(sess)

    // Create an input to get the object from S3
    getObjectInput := &s3.GetObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(key),
    }

    // Get the object from S3
    result, err := s3Client.GetObject(getObjectInput)
    if err != nil {
        log.Printf("Failed to get object from S3: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return err
    }

    // Set the response headers based on the S3 object's metadata
    w.Header().Set("Content-Length", fmt.Sprintf("%d", *result.ContentLength))
    w.Header().Set("Content-Type", *result.ContentType)

    // Copy the S3 object's data to the response body
    _, err = io.Copy(w, result.Body)
    if err != nil {
        log.Printf("Failed to copy S3 object data to response body: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return err
    }

    // Close the response body
    result.Body.Close()
	return nil
}
