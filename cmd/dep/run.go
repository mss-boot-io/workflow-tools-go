/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/4/7 15:54
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/4/7 15:54
 */

package dep

import (
	"github.com/mss-boot-io/workflow-tools/pkg"
	"github.com/mss-boot-io/workflow-tools/pkg/change"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mss-boot-io/workflow-tools/pkg/aws"
	"github.com/mss-boot-io/workflow-tools/pkg/dep"
)

var (
	provider,
	ignorePaths,
	workspace,
	filename,
	projectNameMatch,
	repo,
	mark,
	dependenceMatch string
	storeProvider  string
	bucket, region string

	StartCmd = &cobra.Command{
		Use:          "dep",
		Short:        "exec gradle dependency output leaf service and library",
		Example:      "go-workflow-tools dep",
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
	StartCmd.PersistentFlags().StringVar(&provider,
		"provider", os.Getenv("provider"),
		"code repository provider")
	StartCmd.PersistentFlags().StringVar(&ignorePaths,
		"ignore-paths",
		os.Getenv("ignore_paths"),
		"ignore paths")
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
	StartCmd.PersistentFlags().StringVar(&dependenceMatch,
		"dependence-match",
		os.Getenv("dependence_match"),
		"dependence match")
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
	StartCmd.PersistentFlags().StringVar(&mark,
		"mark", os.Getenv("mark"),
		"commit sha or pull request number")
}

func preRun() {
	if filename == "" {
		filename = "settings.gradle"
	}
	if projectNameMatch == "" {
		projectNameMatch = "rootProject.name =\\s'([^']+)'"
	}
	if dependenceMatch == "" {
		dependenceMatch = "includeBuild\\s'([^']+)'"
	}
}

func run() error {
	services, err := dep.GetAllServices(workspace, filename, projectNameMatch)
	if err != nil {
		log.Println(err)
		return err
	}
	d, err := dep.NewDig(workspace, filename, services, projectNameMatch, dependenceMatch)
	if err != nil {
		log.Println(err)
		return err
	}
	dirs := make([]string, 0)
	//??????change files list
	var files change.Files
	switch storeProvider {
	case "s3":
		err = aws.GetObjectFromS3(region, bucket, change.GetFilename(repo, mark, storeProvider), &files)
	default:
		err = pkg.ReadJsonFile(change.GetFilename("", "", storeProvider), &files)
	}
	if err != nil {
		log.Printf("cmd GetObjectFromS3 error: %s", err.Error())
		return err
	}
	ignores := strings.Split(ignorePaths, ",")
	if ignorePaths == "" {
		ignores = nil
	}
	changeFiles := make([]string, 0)
	changeFiles = append(changeFiles, files.Modified...)
	changeFiles = append(changeFiles, files.Added...)
	changeFiles = append(changeFiles, files.Renamed...)
	for _, c := range changeFiles {
		for _, service := range services {
			if strings.Index(c, strings.Join(service, "/")) > -1 {
				c = strings.Join(service, "/")
				break
			}
		}
		var exist bool
		for i := range ignores {
			if strings.Index(c, ignores[i]) > -1 {
				exist = true
			}
		}
		if !exist {
			dirs = append(dirs, c)
		}
	}
	matrix := d.GetChanged(dirs)

	for i := range matrix {
		matrix[i].ProjectPath, _ = services[matrix[i].Name]
		matrix[i].FindLanguages(workspace)
		if strings.Index(strings.ToLower(matrix[i].Name), dep.Airflow.String()) > -1 {
			matrix[i].Type = dep.Airflow
			continue
		}
		if strings.Index(strings.ToLower(matrix[i].Name), dep.Lambda.String()) > -1 {
			matrix[i].Type = dep.Lambda
			continue
		}
	}
	key := dep.GetFilename(repo, mark, storeProvider)
	//key := fmt.Sprintf("%s/%s/artifact/workflow/service.json", repo, mark)
	switch storeProvider {
	case "s3":
		return aws.PutObjectToS3(region, bucket, key, &matrix, "application/json")
	default:
		//??????????????????
		err = pkg.CreatePath(filepath.Dir(key))
		if err != nil {
			return err
		}
		return pkg.WriteJsonFile(key, &matrix)
	}
}
