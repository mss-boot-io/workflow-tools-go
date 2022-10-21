/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/4/7 11:14
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/4/7 11:14
 */

package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// PutObjectToS3 put object to s3
func PutObjectToS3(region, bucket, key string, data interface{}, contentType string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region))
	if err != nil {
		log.Println(err)
		return err
	}
	client := s3.NewFromConfig(cfg)
	buffer := &bytes.Buffer{}
	switch data.(type) {
	case *bytes.Buffer:
		buffer = data.(*bytes.Buffer)
	case string:
		buffer.WriteString(data.(string))
	case []byte:
		buffer.Write(data.([]byte))
	default:
		err = json.NewEncoder(buffer).Encode(data)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	if contentType == "" {
		contentType = "text/plain"
	}
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        buffer,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetObjectFromS3 get object from s3
func GetObjectFromS3(region, bucket, key string, data interface{}) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region))
	if err != nil {
		log.Println(err)
		return err
	}
	client := s3.NewFromConfig(cfg)
	object, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	defer object.Body.Close()
	err = json.NewDecoder(object.Body).Decode(data)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
