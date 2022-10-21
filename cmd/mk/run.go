/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/4/7 15:23
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/4/7 15:23
 */

package mk

import (
	"encoding/json"
	"fmt"
	"github.com/mss-boot-io/workflow-tools/pkg/cdk8s"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/mss-boot-io/workflow-tools/pkg"
	"github.com/mss-boot-io/workflow-tools/pkg/aws"
	"github.com/mss-boot-io/workflow-tools/pkg/dep"
)

var (
	bucket,
	region,
	repo,
	mark,
	leaf,
	downloadCache,
	uploadCache,
	workspace,
	reportTableFile,
	dockerOrg,
	dockerTags,
	configStage,
	serviceType string
	errorBlock    bool
	dockerPush    bool
	generateCDK8S bool
	storeProvider string
	StartCmd      = &cobra.Command{
		Use:          "mk",
		Short:        "exec  multiple work",
		Example:      "go-workflow-tools mk",
		SilenceUsage: true,
		PreRun: func(_ *cobra.Command, _ []string) {
			log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
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
	StartCmd.PersistentFlags().StringVar(&serviceType,
		"service-type",
		os.Getenv("service_type"),
		"service")
	StartCmd.PersistentFlags().StringVar(&leaf,
		"leaf",
		os.Getenv("leaf"),
		"leaf")
	StartCmd.PersistentFlags().StringVar(&downloadCache,
		"download-cache",
		os.Getenv("download_cache"),
		"download-cache")
	StartCmd.PersistentFlags().StringVar(&uploadCache,
		"upload-cache",
		os.Getenv("upload_cache"),
		"upload-cache")
	StartCmd.PersistentFlags().BoolVar(&errorBlock,
		"error-block",
		false,
		"error-block")
	StartCmd.PersistentFlags().StringVar(&workspace,
		"workspace",
		os.Getenv("workspace"),
		"workspace")
	StartCmd.PersistentFlags().StringVar(&reportTableFile,
		"report-link-file",
		os.Getenv("report-link-file"),
		"report-link-file")
	StartCmd.PersistentFlags().StringVar(&repo,
		"repo", os.Getenv("repo"),
		"repository path(github) or url")
	StartCmd.PersistentFlags().StringVar(&mark,
		"mark", os.Getenv("mark"),
		"commit sha or pull request number")
	StartCmd.PersistentFlags().StringVar(&dockerOrg,
		"docker-org", os.Getenv("docker_organize"),
		"docker org")
	StartCmd.PersistentFlags().BoolVar(&dockerPush,
		"docker-push", false,
		"docker push")
	StartCmd.PersistentFlags().StringVar(&dockerTags,
		"docker-tags", os.Getenv("docker_tags"),
		"docker tags")
	StartCmd.PersistentFlags().BoolVar(&generateCDK8S,
		"generate-cdk8s", false,
		"generate cdk8s")
	StartCmd.PersistentFlags().StringVar(&configStage,
		"config-stage", os.Getenv("config_stage"),
		"config stage")
}

func run() error {
	if !errorBlock {
		switch strings.ToLower(os.Getenv("error_block")) {
		case "t", "true", "1", "y", "yes":
			errorBlock = true
		}
	}
	if serviceType == "" {
		serviceType = dep.Service.String()
	}
	if reportTableFile == "" {
		reportTableFile = ".output/report-link.md"
	}
	leafs := make([]dep.Matrix, 0)

	var err error
	key := dep.GetFilename(repo, mark, storeProvider)
	if leaf != "" {
		err = json.Unmarshal([]byte(leaf), &leafs)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		switch storeProvider {
		case "s3":
			err = aws.GetObjectFromS3(region, bucket, key, &leafs)
		default:
			err = pkg.ReadJsonFile(key, &leafs)
		}
		if err != nil {
			log.Println(err)
			return err
		}
	}

	if len(leafs) > 0 && downloadCache != "" {
		fmt.Printf("######################## %s ########################\n", "Download Cache Starting")
		fmt.Println("#    ", downloadCache)
		err = pkg.Cmd(downloadCache)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Printf("######################## %s ########################\n", "Download Cache Finished")
	}
	for i := range leafs {
		if leafs[i].Type.String() != serviceType {
			continue
		}
		fmt.Printf("######################## %s ########################\n", leafs[i].Name)
		leafs[i].Err = leafs[i].Run(workspace, os.Getenv("cmd"), dockerOrg, dockerTags, dockerPush)
		leafs[i].Finish = true
		fmt.Print("###   ")
		if leafs[i].Err == nil && generateCDK8S {
			fmt.Printf("### generate[%s] %s's cdk8s\n", configStage, leafs[i].Name)
			cdk8s.Generate(filepath.Join(filepath.Join(leafs[i].ProjectPath...), "deploy-config.yml"),
				configStage,
				fmt.Sprintf("%s/%s:%s", dockerOrg, leafs[i].Name, dockerTags))

		}
		if leafs[i].Err != nil {
			fmt.Println("Failed")
		} else {
			fmt.Println("Successful")
		}
		fmt.Printf("######################## %s ########################\n", leafs[i].Name)
		if leafs[i].Err != nil && errorBlock {
			break
		}
	}

	fmt.Println()
	fmt.Println()
	fmt.Println()

	var failed bool
	defer func() {
		if failed {
			os.Exit(-1)
		}
	}()
	fmt.Printf("######################## %s ########################\n", "All Service")
	var updateCache bool
	for i := range leafs {
		if leafs[i].Type == dep.Library {
			for j := range leafs[i].Language {
				if leafs[i].Language[j] == dep.JAVA {
					updateCache = true
					break
				}
			}
		}
		if !leafs[i].Finish {
			continue
		}
		if leafs[i].Type != dep.Service {
			continue
		}
		fmt.Printf("###   %s[%s]: ", leafs[i].Name, leafs[i].LanguageString())
		if leafs[i].Err != nil {
			failed = true
			err = leafs[i].Err
			fmt.Println("Failed")
			continue
		}
		fmt.Println("Successful")
	}
	fmt.Printf("######################## %s ########################\n", "All Service")
	// Generate test report link table
	_ = dep.OutputReportTableToPR(os.Getenv("repo_path"), cast.ToInt(os.Getenv("pr_number")), leafs)
	if failed {
		return err
	}
	if err != nil {
		log.Println(err)
		return err
	}
	if len(leafs) > 0 && uploadCache != "" && updateCache {
		fmt.Printf("######################## %s ########################\n", "Upload Cache Starting")
		fmt.Println("#    ", uploadCache)
		err = pkg.Cmd(uploadCache)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Printf("######################## %s ########################\n", "Upload Cache Finished")
	}
	return err
}
