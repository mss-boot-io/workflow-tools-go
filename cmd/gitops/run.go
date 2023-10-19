/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/11/1 06:52:05
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/11/1 06:52:05
 */

package gitops

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mss-boot-io/workflow-tools/pkg"
	"github.com/mss-boot-io/workflow-tools/pkg/argocd"
	"github.com/mss-boot-io/workflow-tools/pkg/aws"
	"github.com/mss-boot-io/workflow-tools/pkg/dep"
	"github.com/mss-boot-io/workflow-tools/pkg/gitops"
)

var (
	repo,
	mark,
	bucket,
	region,
	leaf,
	configStage,
	argocdURL,
	argocdToken,
	argocdProject,
	argocdNamespace,
	gitopsRepo,
	gitopsBranch,
	gitopsConfigFile,
	workspace,
	storeProvider string
	languageEnv  string
	singleGitops string
	errorBlock   bool
	StartCmd     = &cobra.Command{
		Use:          "gitops",
		Short:        "exec  multiple work to gitops",
		Example:      "go-workflow-tools gitops",
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
	StartCmd.PersistentFlags().StringVar(&workspace,
		"workspace",
		os.Getenv("workspace"),
		"workspace")
	StartCmd.PersistentFlags().StringVar(&gitopsConfigFile,
		"gitops-config-file",
		"deploy-config.yml",
		"gitops config file name")
	StartCmd.PersistentFlags().StringVar(&argocdURL,
		"argocd-url",
		os.Getenv("argocd_url"),
		"argocd url")
	StartCmd.PersistentFlags().StringVar(&argocdToken,
		"argocd-token",
		os.Getenv("argocd_token"),
		"argocd token")
	StartCmd.PersistentFlags().StringVar(&argocdProject,
		"argocd-project",
		os.Getenv("argocd_project"),
		"argocd project")
	StartCmd.PersistentFlags().StringVar(&argocdNamespace,
		"argocd-namespace",
		"argocd",
		"argocd namespace")
	StartCmd.PersistentFlags().StringVar(&gitopsRepo,
		"gitops-repo",
		os.Getenv("gitops_repo"),
		"gitops repo")
	StartCmd.PersistentFlags().StringVar(&gitopsBranch,
		"gitops-branch",
		os.Getenv("gitops_branch"),
		"gitops branch")
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
	StartCmd.PersistentFlags().StringVar(&leaf,
		"leaf",
		os.Getenv("leaf"),
		"leaf")
	StartCmd.PersistentFlags().StringVar(&repo,
		"repo", os.Getenv("repo"),
		"repository path(github) or url")
	StartCmd.PersistentFlags().StringVar(&mark,
		"mark", os.Getenv("mark"),
		"commit sha or pull request number")
	StartCmd.PersistentFlags().StringVar(&configStage,
		"config-stage", os.Getenv("config_stage"),
		"config stage")
	StartCmd.PersistentFlags().BoolVar(&errorBlock,
		"error-block",
		false,
		"error-block")
	StartCmd.PersistentFlags().StringVar(&languageEnv,
		"languageEnv",
		os.Getenv("languageEnv"),
		"the language and version required for the build")
	StartCmd.PersistentFlags().StringVar(&singleGitops,
		"singleGitops",
		os.Getenv("singleGitops"),
		"only supported build languages and corresponding versions")
}

func preRun() {
	if singleGitops == "" {
		singleGitops = "false"
	}
}

func run() error {
	if !errorBlock {
		errorBlock = cast.ToBool(os.Getenv("error_block"))
	}
	var failed bool
	defer func() {
		if failed {
			os.Exit(-1)
		}
	}()
	serviceType := dep.Service.String()
	leafs := make([]dep.Matrix, 0)
	var err error
	key := dep.GetFilename(repo, mark, storeProvider)
	if leaf != "" && leaf != "[]" {
		err = json.Unmarshal([]byte(leaf), &leafs)
		if err != nil {
			fmt.Printf("unmarshal leaf error: %s\n", err.Error())
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
			fmt.Printf("get leafs from %s error: %s\n", key, err.Error())
			return err
		}
	}
	// When only a single gitops is supported, parse out the supported languages and versions.
	var languageType, languageVersion string
	isSingleGitops, err := strconv.ParseBool(singleGitops)
	if err != nil {
		return err
	}
	if isSingleGitops {
		pattern := "/"
		match, err := regexp.MatchString(pattern, languageEnv)
		if err != nil || !match {
			return err
		}
		languageType = strings.Split(languageEnv, "/")[0]
		languageVersion = strings.Split(languageEnv, "/")[1]
	}
	for i := range leafs {
		if leafs[i].Type.String() != serviceType {
			continue
		}
		if isSingleGitops && (leafs[i].LanguageEnvType != languageType || leafs[i].LanguageEnvVersion != languageVersion) {
			continue
		}
		fmt.Printf("######################## %s ########################\n", leafs[i].Name)
		if len(leafs[i].ProjectPath) == 0 {
			leafs[i].ProjectPath = []string{leafs[i].Name}
		}
		var gitopsConfig *gitops.Config
		fmt.Print("###   \n")
		gitopsConfig, leafs[i].Err = gitops.LoadFile(filepath.Join(workspace, filepath.Join(leafs[i].ProjectPath...), gitopsConfigFile))
		if leafs[i].Err != nil {
			fmt.Printf("### load %s's gitops config file error: %s\n", leafs[i].Name, leafs[i].Err)
			continue
		}
		if argocdURL == "" || argocdToken == "" {
			failed = true
			err = errors.New("argocd url or token is empty")
			return err
		}
		argocdClient := argocd.New(argocdURL, argocdToken, nil)
		for stage := range gitopsConfig.Deploy.Stage {
			if (strings.Index(stage, "prod") > -1 || strings.Index(stage, "production") > -1) &&
				!(configStage == "prod" || configStage == "production") {
				continue
			}
			if (strings.Index(stage, "prod") == -1 && strings.Index(stage, "production") == -1) &&
				(configStage == "prod" || configStage == "production") {
				continue
			}
			namespace := gitopsConfig.Deploy.Stage[stage].GetKey("namespace")
			if namespace == "" {
				namespace = stage
			}
			paths := make([]string, 0)
			if gitopsConfig.Project != "" && strings.Index(leafs[i].Name, gitopsConfig.Project) == -1 {
				paths = append(paths, gitopsConfig.Project)
			}
			paths = append(paths, leafs[i].Name, stage)
			app := &appv1.Application{
				ObjectMeta: metav1.ObjectMeta{
					Name:      strings.Join(paths, "-"),
					Namespace: argocdNamespace,
					Labels: map[string]string{
						"Project": gitopsConfig.Project,
					},
				},
				Spec: appv1.ApplicationSpec{
					Project: argocdProject,
					Source: &appv1.ApplicationSource{
						RepoURL:        gitopsRepo,
						Path:           fmt.Sprintf("%s/%s", configStage, strings.Join(leafs[i].ProjectPath, "/")),
						TargetRevision: gitopsBranch,
					},
					Destination: appv1.ApplicationDestination{
						Name:      cast.ToString(gitopsConfig.Deploy.Stage[stage].GetKey("cluster")),
						Namespace: cast.ToString(namespace),
					},
				},
			}
			if cast.ToBool(gitopsConfig.Deploy.Stage[stage].GetKey("autoSync")) {
				app.Spec.SyncPolicy = &appv1.SyncPolicy{
					Automated: &appv1.SyncPolicyAutomated{
						Prune:      true,
						SelfHeal:   true,
						AllowEmpty: true,
					},
				}
			}
			err = argocdClient.CreateApplication(app)
			if err != nil {
				fmt.Printf("### create %s's %s stage application error: %s\n", leafs[i].Name, stage, err)
				leafs[i].Err = err
				break
			}
			fmt.Printf("### create %s's %s stage application success\n", leafs[i].Name, stage)
		}
		fmt.Printf("### gitops config successed\n")
		fmt.Printf("######################## %s ########################\n", leafs[i].Name)
		if leafs[i].Err != nil && errorBlock {
			break
		}
	}

	for i := range leafs {
		if leafs[i].Err != nil {
			failed = true
			err = leafs[i].Err
		}
	}
	return err
}
