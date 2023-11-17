/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/8/29 10:39:38
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/8/29 10:39:38
 */

package change

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mss-boot-io/workflow-tools/pkg"
	"github.com/mss-boot-io/workflow-tools/pkg/aws"
	"github.com/mss-boot-io/workflow-tools/pkg/change"
	"github.com/mss-boot-io/workflow-tools/pkg/change/github"
	"github.com/mss-boot-io/workflow-tools/pkg/minio"
)

var (
	provider,
	accessToken,
	mark,
	storeProvider,
	repo string
	minioEndpoint        string
	minioAccessKey       string
	minioSecretAccessKey string
	bucket, region       string
	StartCmd             = &cobra.Command{
		Use:          "change",
		Short:        "exec change for the repo",
		Example:      "go-workflow-tools change",
		SilenceUsage: true,
		PreRun: func(_ *cobra.Command, _ []string) {
			log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
			preRun()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVar(&storeProvider,
		"store-provider",
		os.Getenv("storeProvider"),
		"store provider")
	StartCmd.PersistentFlags().StringVar(&provider,
		"provider", os.Getenv("provider"),
		"code repository provider")
	StartCmd.PersistentFlags().StringVar(&accessToken,
		"accessToken", os.Getenv("accessToken"),
		"code repository accessToken")
	StartCmd.PersistentFlags().StringVar(&repo,
		"repo", os.Getenv("repo"),
		"repository path(github) or url")
	StartCmd.PersistentFlags().StringVar(&mark,
		"mark", os.Getenv("mark"),
		"commit sha or pull request number")
	StartCmd.PersistentFlags().StringVar(&bucket,
		"bucket",
		os.Getenv("bucket"),
		"bucket")
	StartCmd.PersistentFlags().StringVar(&region,
		"region",
		os.Getenv("region"),
		"region")
	StartCmd.PersistentFlags().StringVar(&minioEndpoint,
		"minioEndpoint", os.Getenv("minioEndpoint"),
		"minioEndpoint")
	StartCmd.PersistentFlags().StringVar(&minioAccessKey,
		"minioAccessKey", os.Getenv("minioAccessKey"),
		"minioAccessKey")
	StartCmd.PersistentFlags().StringVar(&minioSecretAccessKey,
		"minioSecretAccessKey", os.Getenv("minioSecretAccessKey"),
		"minioSecretAccessKey")
}

func preRun() {
	if storeProvider == "" {
		storeProvider = "s3"
	}
}

func run() error {
	switch strings.ToLower(provider) {
	case "github", "":
		o := github.Github{}
		o.SetAuth(accessToken)
		err := o.SetRepoURL(repo)
		if err != nil {
			log.Printf("cmd SetRepoURL error: %s", err.Error())
			return err
		}

		files, err := o.ChangeFiles(mark)
		if err != nil {
			log.Printf("cmd ChangeFiles error: %s", err.Error())
			return err
		}
		key := change.GetFilename(repo, mark, storeProvider)
		switch storeProvider {
		case "s3":
			return aws.PutObjectToS3(region, bucket, key, *files, "")
		case "minio":
			minioCli := minio.New(minioEndpoint, minioAccessKey, minioSecretAccessKey)
			return minioCli.PutObject(bucket, key, *files)
		default:
			//默认使用文件
			err = pkg.CreatePath(filepath.Dir(key))
			if err != nil {
				return err
			}
			return pkg.WriteJsonFile(key, files)
		}
	default:
		log.Fatalf("not support %s", provider)
	}
	return nil
}
