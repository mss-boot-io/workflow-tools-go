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
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/mss-boot-io/workflow-tools/pkg"
	"github.com/mss-boot-io/workflow-tools/pkg/aws"
	"github.com/mss-boot-io/workflow-tools/pkg/cdk8s"
	"github.com/mss-boot-io/workflow-tools/pkg/dep"
	"github.com/mss-boot-io/workflow-tools/pkg/gitops"
	"github.com/mss-boot-io/workflow-tools/pkg/minio"
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
	dockerTags,
	configStage,
	gitopsConfigFile,
	rulesFilename,
	alertManagerPrefix,
	alertManagerUsername,
	alertManagerPassword,
	serviceType string
	errorBlock           bool
	dockerPush           bool
	generateCDK8S        bool
	storeProvider        string
	makefileTmplPath     string
	languageEnv          string
	singleBuildEnv       string
	minioEndpoint        string
	minioAccessKey       string
	minioSecretAccessKey string
	minioSSL             string
	StartCmd             = &cobra.Command{
		Use:          "mk",
		Short:        "exec  multiple work",
		Example:      "go-workflow-tools mk",
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
	StartCmd.PersistentFlags().StringVar(&gitopsConfigFile,
		"gitops-config-file",
		"deploy-config.yml",
		"gitops config file name")
	StartCmd.PersistentFlags().StringVar(&storeProvider,
		"store-provider",
		os.Getenv("storeProvider"),
		"store provider")
	StartCmd.PersistentFlags().StringVar(&bucket,
		"bucket",
		os.Getenv("bucket"),
		"bucket")
	StartCmd.PersistentFlags().StringVar(&region,
		"region",
		os.Getenv("region"),
		"region")
	StartCmd.PersistentFlags().StringVar(&rulesFilename,
		"rules-filename",
		"alert_manager_rules",
		"rules-filename")
	StartCmd.PersistentFlags().StringVar(&alertManagerPrefix,
		"alert-manager-prefix",
		os.Getenv("alert_manager_prefix"),
		"alert-manager-prefix")
	StartCmd.PersistentFlags().StringVar(&alertManagerUsername,
		"alert-manager-username",
		os.Getenv("alert_manager_username"),
		"alert-manager-username")
	StartCmd.PersistentFlags().StringVar(&alertManagerPassword,
		"alert-manager-password",
		os.Getenv("alert_manager_password"),
		"alert-manager-password")
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
	StartCmd.PersistentFlags().StringVar(&makefileTmplPath,
		"makefileTmplPath",
		os.Getenv("makefileTmplPath"),
		"makefile template path")
	StartCmd.PersistentFlags().StringVar(&languageEnv,
		"languageEnv",
		os.Getenv("languageEnv"),
		"the language and version required for the build")
	StartCmd.PersistentFlags().StringVar(&singleBuildEnv,
		"singleBuildEnv",
		os.Getenv("singleBuildEnv"),
		"Only supports a single build locale")
	StartCmd.PersistentFlags().StringVar(&minioEndpoint,
		"minioEndpoint", os.Getenv("minioEndpoint"),
		"minioEndpoint")
	StartCmd.PersistentFlags().StringVar(&minioAccessKey,
		"minioAccessKey", os.Getenv("minioAccessKey"),
		"minioAccessKey")
	StartCmd.PersistentFlags().StringVar(&minioSecretAccessKey,
		"minioSecretAccessKey", os.Getenv("minioSecretAccessKey"),
		"minioSecretAccessKey")
	StartCmd.PersistentFlags().StringVar(&minioSSL,
		"minioSSL", os.Getenv("minioSSL"),
		"use ssl in minio transmission")
}

func preRun() {
	if singleBuildEnv == "" {
		singleBuildEnv = "false"
	}
	if storeProvider == "" {
		storeProvider = "s3"
	}
	if minioSSL == "" {
		minioSSL = "true"
	}
}

func run() error {
	if !errorBlock {
		errorBlock = cast.ToBool(os.Getenv("error_block"))
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
	if leaf != "" && leaf != "[]" {
		err = json.Unmarshal([]byte(leaf), &leafs)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		switch storeProvider {
		case "s3":
			err = aws.GetObjectFromS3(region, bucket, key, &leafs)
		case "minio":
			minioCli := minio.New(minioEndpoint, minioAccessKey, minioSecretAccessKey, minioSSL)
			err = minioCli.GetObject(bucket, key, &leafs)
		default:
			err = pkg.ReadJsonFile(key, &leafs)
		}
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// When only a single build is supported, parse out the supported languages and versions.
	var languageType, languageVersion string
	isSingleBuildoEnv, err := strconv.ParseBool(singleBuildEnv)
	if err != nil {
		return err
	}
	if isSingleBuildoEnv {
		pattern := "/"
		match, err := regexp.MatchString(pattern, languageEnv)
		if err != nil || !match {
			return err
		}
		languageType = strings.Split(languageEnv, "/")[0]
		languageVersion = strings.Split(languageEnv, "/")[1]
	}

	// Download cache if needed for any project
	var saveCacheNeeded bool
	for i := range leafs {
		if isSingleBuildoEnv && (leafs[i].LanguageEnvType != languageType || leafs[i].LanguageEnvVersion != languageVersion) {
			continue
		}
		if leafs[i].ProjectPath == nil || len(leafs[i].ProjectPath) == 0 {
			leafs[i].ProjectPath = []string{leafs[i].Name}
		}
		if strings.Index(serviceType, leafs[i].Type.String()) == -1 {
			continue
		}
		var config *gitops.Config
		config, leafs[i].Err = gitops.LoadFile(filepath.Join(workspace, filepath.Join(leafs[i].ProjectPath...), gitopsConfigFile))
		if leafs[i].Err != nil && errorBlock {
			break
		}
		if config != nil && !config.Build.SkipCache {
			saveCacheNeeded = true
			break
		}
	}
	if saveCacheNeeded && downloadCache != "" {
		fmt.Printf("######################## %s ########################\n", "Download Cache Starting")
		fmt.Println("#    ", downloadCache)
		downloadErr := pkg.Cmd(downloadCache)
		if downloadErr != nil {
			log.Println(downloadErr)
		}
		fmt.Printf("######################## %s ########################\n", "Download Cache Finished")
	}

	for i := range leafs {
		if strings.Index(serviceType, leafs[i].Type.String()) == -1 {
			continue
		}

		if isSingleBuildoEnv && (leafs[i].LanguageEnvType != languageType || leafs[i].LanguageEnvVersion != languageVersion) {
			continue
		}

		fmt.Printf("######################## %s ########################\n", leafs[i].Name)
		var gitopsConfig *gitops.Config
		gitopsConfig, leafs[i].Err = gitops.LoadFile(filepath.Join(workspace, filepath.Join(leafs[i].ProjectPath...), gitopsConfigFile))
		if leafs[i].Err != nil && errorBlock {
			break
		}
		leafs[i].Project = gitopsConfig.Project
		var dockerImage string
		if gitopsConfig != nil {
			dockerImage = gitopsConfig.GetImage(leafs[i].Name)
		}

		var cmd string
		if leafs[i].Type != dep.Service && os.Getenv(fmt.Sprintf("%s_cmd", leafs[i].Type.String())) != "" {
			// get service type cmd
			cmd = os.Getenv(fmt.Sprintf("%s_cmd", leafs[i].Type.String()))
			var stage string
			switch configStage {
			case "prod":
				stage = "prod"
			default:
				stage = "alpha"
			}
			stageDeploy, ok := gitopsConfig.Deploy.Stage[stage]
			if !ok {
				leafs[i].Err = fmt.Errorf("%s stage not found", stage)
				if errorBlock {
					break
				}
			}
			cmd, leafs[i].Err = stageDeploy.ParseTemplate(cmd)
			if leafs[i].Err != nil && errorBlock {
				break
			}

		}
		if cmd != "" {
			cmd = "&" + cmd
		}
		cmd = os.Getenv("cmd") + cmd
		fmt.Printf("### cmd: %s\n", cmd)

		// copy makefile template to service path if not exist
		makefilePath := filepath.Join(workspace, filepath.Join(leafs[i].ProjectPath...), "Makefile")
		makefileExist := pkg.PathExist(makefilePath)
		if !makefileExist {
			if makefileTmplPath == "" {
				makefileTmplPath = filepath.Join(workspace, ".github/Makefile_Temp")
			}
			err := pkg.CopyFile(makefileTmplPath, makefilePath)
			if err != nil {
				log.Println(err)
				return err
			}
		}

		leafs[i].Err = leafs[i].Run(workspace, os.Getenv("cmd"), dockerImage, dockerTags, dockerPush)
		leafs[i].Finish = true
		fmt.Print("###   ")
		if leafs[i].Type == dep.Service && leafs[i].Err == nil && generateCDK8S {
			fmt.Printf("### generate[%s] %s's cdk8s\n", configStage, leafs[i].Name)
			cdk8s.Generate(filepath.Join(filepath.Join(leafs[i].ProjectPath...), "deploy-config.yml"),
				configStage,
				fmt.Sprintf("%s:%s", dockerImage, dockerTags),
				leafs[i].ProjectPath)

		}
		err = leafs[i].Alert(workspace, rulesFilename, alertManagerPrefix, alertManagerUsername, alertManagerPassword)
		if err != nil {
			// ignore alert error
			fmt.Printf("### ignore alert rule error: %s\n", err.Error())
			err = nil
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
		if isSingleBuildoEnv && (leafs[i].LanguageEnvType != languageType || leafs[i].LanguageEnvVersion != languageVersion) {
			continue
		}
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

	// Upload cache if needed for any project
	if saveCacheNeeded && uploadCache != "" && updateCache {
		fmt.Printf("######################## %s ########################\n", "Upload Cache Starting")
		fmt.Println("#    ", uploadCache)
		uploadErr := pkg.Cmd(uploadCache)
		if uploadErr != nil {
			log.Println(uploadErr)
		}
		fmt.Printf("######################## %s ########################\n", "Upload Cache Finished")
	}
	return err
}
