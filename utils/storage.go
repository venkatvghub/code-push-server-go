package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage interface {
	SaveFile(filePath, key string) error
	GetFileURL(key string) string
	UploadFile(filePath, key string) error
}

func NewStorage() Storage {
	switch Config.Storage.Type {
	case "s3":
		return NewS3Storage()
	default: // "local" or unrecognized falls back to local
		return NewLocalStorage()
	}
}

// LocalStorage implementation
type LocalStorage struct{}

func NewLocalStorage() Storage {
	return &LocalStorage{}
}

func (s *LocalStorage) SaveFile(filePath, dest string) error {
	err := os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}
	return os.Rename(filePath, dest)
}

func (s *LocalStorage) GetFileURL(fileName string) string {
	return Config.Storage.Local.DownloadUrl + "/" + filepath.Base(fileName)
}

func (s *LocalStorage) UploadFile(filePath, key string) error {
	return s.SaveFile(filePath, Config.Storage.Local.StorageDir+"/"+key)
}

// S3Storage implementation
type S3Storage struct {
	client *s3.Client
}

func NewS3Storage() Storage {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(Config.Storage.S3.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			Config.Storage.S3.AccessKeyID,
			Config.Storage.S3.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		panic("Failed to load AWS config: " + err.Error())
	}

	return &S3Storage{
		client: s3.NewFromConfig(cfg),
	}
}

func (s *S3Storage) SaveFile(filePath, key string) error {
	return s.UploadFile(filePath, key)
}

func (s *S3Storage) GetFileURL(key string) string {
	if Config.Storage.S3.DownloadUrl != "" {
		return Config.Storage.S3.DownloadUrl + "/" + key
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", Config.Storage.S3.BucketName, Config.Storage.S3.Region, key)
}

func (s *S3Storage) UploadFile(filePath, key string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(Config.Storage.S3.BucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	return err
}
