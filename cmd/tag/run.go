/*
 * @Author: snakelu<lyhccq@163.com>
 * @Date: 2023/5/31 09:18
 * @Last Modified by: snakelu<lyhccq@163.com>
 * @Last Modified time: 2023/5/31 09:18
 */

package tag

import (
	"errors"
	"github.com/mss-boot-io/workflow-tools/pkg"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mss-boot-io/workflow-tools/pkg/aws"
	"github.com/mss-boot-io/workflow-tools/pkg/dep"
)

var (
	ref,
	provider,
	workspace,
	filename,
	projectNameMatch,
	repo,
	storeProvider string
	bucket, region string

	StartCmd = &cobra.Command{
		Use:          "tag",
		Short:        "exec gradle dependency output leaf service and library",
		Example:      "go-workflow-tools tag",
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
	StartCmd.PersistentFlags().StringVar(&ref,
		"ref", os.Getenv("ref"),
		"code repository tag ref")
	StartCmd.PersistentFlags().StringVar(&provider,
		"provider", os.Getenv("provider"),
		"code repository provider")
	StartCmd.PersistentFlags().StringVar(&workspace,
		"workspace",
		os.Getenv("workspace"),
		"workspace")
	StartCmd.PersistentFlags().StringVar(&filename,
		"filename",
		os.Getenv("filename"),
		"filename")
	StartCmd.PersistentFlags().StringVar(&projectNameMatch,
		"project-name-match",
		os.Getenv("project_name_match"),
		"project name match")
	StartCmd.PersistentFlags().StringVar(&storeProvider,
		"store-provider",
		"s3",
		"store provider")
	StartCmd.PersistentFlags().StringVar(&bucket,
		"bucket",
		os.Getenv("bucket"),
		"bucket")
	StartCmd.PersistentFlags().StringVar(&region,
		"region",
		os.Getenv("region"),
		"region")
	StartCmd.PersistentFlags().StringVar(&repo,
		"repo", os.Getenv("repo"),
		"repository path(github) or url")
}

func preRun() {
	if filename == "" {
		filename = "settings.gradle"
	}
	if projectNameMatch == "" {
		projectNameMatch = "rootProject.name =\\s'([^']+)'"
	}
}

func run() error {
	// ref format example: refs/tags/terminal/v0.0.1
	if !strings.HasPrefix(ref, "refs/tags/") {
		log.Printf("ref is invalid %s\n", ref)
		return errors.New("ref is invalid")
	}

	serviceName := strings.Split(ref, "/")[2]
	services, err := dep.GetAllServices(workspace, filename, projectNameMatch)
	if err != nil {
		log.Println(err)
		return err
	}

	if services[serviceName] == nil {
		log.Printf("service %s not exist\n", serviceName)
		return errors.New("service not exist")
	}

	matrix := &dep.Matrix{
		Name:        serviceName,
		Type:        dep.Service,
		ProjectPath: services[serviceName],
	}

	matrix.FindLanguages(workspace)
	if strings.Index(strings.ToLower(matrix.Name), dep.Airflow.String()) > -1 {
		matrix.Type = dep.Airflow
	}
	if strings.Index(strings.ToLower(matrix.Name), dep.Lambda.String()) > -1 {
		matrix.Type = dep.Lambda
	}

	matrices := make([]*dep.Matrix, 0, 1)
	matrices[0] = matrix

	key := dep.GetFilename(repo, ref, storeProvider)
	//key := fmt.Sprintf("%s/%s/artifact/workflow/service.json", repo, mark)
	switch storeProvider {
	case "s3":
		return aws.PutObjectToS3(region, bucket, key, &matrices, "application/json")
	default:
		//默认使用文件
		err = pkg.CreatePath(filepath.Dir(key))
		if err != nil {
			return err
		}
		return pkg.WriteJsonFile(key, &matrix)
	}
}
