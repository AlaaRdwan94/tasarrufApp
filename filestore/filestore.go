package filestore

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

// File is the file to be pushed to S3 bucket
type File struct {
	Body        []byte
	FileName    string
	Size        int64
	ContentType string
}

// UploadToS3 uploads the given file to amazon s3 bucket
func (f *File) UploadToS3() (*s3manager.UploadOutput, error) {
	region := os.Getenv("S3_REGION")
	bucketName := os.Getenv("BUCKET_NAME")
	Akid := os.Getenv("AWS_ACCESS_KEY")
	Secretkey := os.Getenv("AWS_SECRET_KEY")
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(Akid, Secretkey,""),
	})
	if err != nil {
		log.Error(err)
	}
	object := &s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(f.FileName),
		ACL:    aws.String("public-read"),
		Body:   bytes.NewReader(f.Body),
	}
	uploader := s3manager.NewUploader(s)
	out, err := uploader.Upload(object)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return out, nil
}

// DeleteFromS3 deletes a file from s3 bucket.
func DeleteFromS3(fileKey string) error {
	region := os.Getenv("S3_REGION")
	bucketname := os.Getenv("BUCKET_NAME")
	if fileKey == "default.png" {
		return nil
	}
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return err
	}
	svc := s3.New(s)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketname), Key: aws.String(fileKey)})
	if err != nil {
		return err
	}
	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return err
	}
	return nil
}
