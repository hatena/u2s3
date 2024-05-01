package aws

import (
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/hatena/u2s3/pkg/config"
)

type S3Cli struct {
	s3Svc  s3iface.S3API
	bucket string
}

func NewS3Cli(config *config.UploadConfig) *S3Cli {
	if os.Getenv("CUSTOM_HOST") != "" {
		return NewCliForCustom(config)
	}
	sess, err := session.NewSession()
	sess.Config.MaxRetries = aws.Int(config.MaxRetry)
	if err != nil {
		log.Fatal("Creating session is failed")
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		log.Fatal("You must specify a region")
	}
	s3Svc := s3.New(sess)
	return &S3Cli{s3Svc, config.Bucket}
}

func NewCliForCustom(config *config.UploadConfig) *S3Cli {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), ""),
		Endpoint:         aws.String(os.Getenv("CUSTOM_HOST")),
		Region:           aws.String(os.Getenv("CUSTOM_REGION")),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		MaxRetries:       aws.Int(config.MaxRetry),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatal("Creating session is failed")
	}
	s3Svc := s3.New(newSession)
	return &S3Cli{s3Svc, config.Bucket}
}

func (c *S3Cli) Upload(key string, body io.ReadSeeker) error {
	uploader := s3manager.NewUploaderWithClient(c.s3Svc)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	return err
}

func (c *S3Cli) HasKey(key string) bool {
	_, err := c.s3Svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	return err == nil
}
