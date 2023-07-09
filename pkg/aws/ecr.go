package aws

/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2023/7/9 15:15:07
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2023/7/9 15:15:07
 */

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func CreatePrivateRepoIfNotExist(region, image string) error {
	index := strings.Index(image, "/")
	if index == -1 {
		return fmt.Errorf("invalid image: %s", image)
	}
	image = image[index+1:]
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Printf("load config error: %v\n", err)
		return err
	}
	client := ecr.NewFromConfig(cfg)
	// check image exist
	_, err = client.DescribeRepositories(context.TODO(), &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{image},
	})
	if err == nil {
		return nil
	}
	// create image
	_, err = client.CreateRepository(context.TODO(), &ecr.CreateRepositoryInput{
		RepositoryName: &image,
	})
	if err != nil {
		log.Printf("create repository error: %v\n", err)
		return err
	}
	return nil
}
