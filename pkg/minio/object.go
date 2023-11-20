package minio

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client *minio.Client
}

func New(endpoint, accessKey, secretAccessKey, useSSL string) *MinioClient {
	minioClient := &MinioClient{}
	isSSL, _ := strconv.ParseBool(useSSL)
	minioClient.client, _ = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
		Secure: isSSL,
	})
	return minioClient
}

func (c MinioClient) PutObject(bucketName, objectName string, data interface{}) error {
	buffer := &bytes.Buffer{}
	switch data.(type) {
	case *bytes.Buffer:
		buffer = data.(*bytes.Buffer)
	case string:
		buffer.WriteString(data.(string))
	case []byte:
		buffer.Write(data.([]byte))
	default:
		err := json.NewEncoder(buffer).Encode(data)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	_, err := c.client.PutObject(context.Background(), bucketName, objectName, buffer, int64(buffer.Len()), minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c MinioClient) GetObject(bucketName, objectName string, data interface{}) error {
	object, err := c.client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()

	err = json.NewDecoder(object).Decode(data)
	if err != nil {
		return err
	}
	return nil
}
