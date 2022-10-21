/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/4/6 11:19
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/4/6 11:19
 */

package dep

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v44/github"
	"github.com/zhnxin/csvreader"
	"golang.org/x/oauth2"
)

type (
	ServiceType string
	Language    string
)

const (
	Service ServiceType = "service"
	Library ServiceType = "library"
	Lambda  ServiceType = "lambda"
	Airflow ServiceType = "airflow"

	Golang Language = "golang"
	JAVA   Language = "java"
	Python Language = "python"
	Node   Language = "node"
)

type Matrix struct {
	Name        string      `json:"name"`
	Type        ServiceType `json:"type"`
	Language    []Language  `json:"language"`
	ReportUrl   string      `json:"report"`
	Err         error       `json:"-"`
	ProjectPath []string    `json:"projectPath"`
	Finish      bool        `json:"-"`
	Reports     []Report    `json:"_"`
	Coverage    float64     `json:"coverage"`
}

func (e *Matrix) LanguageString() string {
	if len(e.Language) == 0 {
		return ""
	}
	languages := make([]string, len(e.Language))
	for i := range e.Language {
		languages[i] = e.Language[i].String()
	}
	return strings.Join(languages, ",")
}

func (e *Matrix) Run(workspace, cs, dockerOrg, dockerTags string, dockerPush bool) error {
	if len(e.ProjectPath) == 0 {
		e.ProjectPath = []string{e.Name}
	}
	if workspace != "" {
		cs = fmt.Sprintf("cd %s && %s", filepath.Join(workspace, filepath.Join(e.ProjectPath...)), cs)
	} else {
		cs = fmt.Sprintf("cd %s && %s", filepath.Join(e.ProjectPath...), cs)
	}
	if dockerOrg != "" {
		cs += fmt.Sprintf(" && docker build -t %s/%s:latest .", dockerOrg, e.Name)
	}
	if dockerTags != "" {
		var pushLatest bool
		for _, tag := range strings.Split(dockerTags, ",") {
			if strings.Index(tag, "v") > -1 && len(tag) < 10 {
				pushLatest = true
			}
			cs += fmt.Sprintf(" && docker tag %s/%s:latest %s/%s:%s", dockerOrg, e.Name, dockerOrg, e.Name, tag)
			if dockerPush {
				cs += fmt.Sprintf(" && docker push %s/%s:%s", dockerOrg, e.Name, tag)
			}
		}
		if pushLatest && dockerPush {
			cs += fmt.Sprintf(" && docker push %s/%s:latest", dockerOrg, e.Name)
		}
	}
	cmd := exec.Command("/bin/bash", "-c", cs)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s ServiceType) String() string {
	return string(s)
}

func (l Language) String() string {
	return string(l)
}

func (e *Matrix) FindLanguages(workspace string) {
	var (
		golang, python, java, node bool
	)
	_ = filepath.WalkDir(filepath.Join(workspace, filepath.Join(e.ProjectPath...)), func(path string, d fs.DirEntry, err error) error {
		switch filepath.Ext(path) {
		case ".go":
			golang = true
		case ".java":
			java = true
		case ".py":
			python = true
		case ".ts", ".js", ".tsx", ".jsx":
			node = true
		}
		return nil
	})
	if golang {
		e.Language = append(e.Language, Golang)
	}
	if python {
		e.Language = append(e.Language, Python)
	}
	if java {
		e.Language = append(e.Language, JAVA)
	}
	if node {
		e.Language = append(e.Language, Node)
	}
}

// OutputReportTableToPR output report table to PR
func OutputReportTableToPR(repoPath string, number int, list []Matrix) (err error) {
	if repoPath == "" || len(strings.Split(repoPath, "/")) < 2 {
		return fmt.Errorf("repoPath is invalid")
	}
	if number <= 0 {
		return fmt.Errorf("number is invalid")
	}
	var markdownStr string
	coverageStandard := cast.ToFloat64(os.Getenv("coverage_standard"))
	for i := range list {
		//support java
		if list[i].ReportUrl == "" ||
			strings.Index(list[i].LanguageString(), "java") == -1 {
			continue
		}
		list[i].Reports = make([]Report, 0)
		_ = csvreader.New().UnMarshalFile(
			filepath.Join(filepath.Join(list[i].ProjectPath...),
				"build",
				"reports",
				"jacoco",
				"coverage.csv"),
			&(list[i].Reports))
		//Calculate coverage
		var covered, missed int
		for j := range list[i].Reports {
			covered += list[i].Reports[j].InstructionCovered
			missed += list[i].Reports[j].InstructionMissed
		}
		if covered+missed > 0 {
			list[i].Coverage = float64(covered) / float64(covered+missed) * 100
		}
		emoji := ":green_apple:"
		if list[i].Coverage < coverageStandard {
			emoji = ":x:"
		}

		//java service
		markdownStr += fmt.Sprintf("| %s | %s | %s | %s %.2f%s | [go to report](%s) |\n",
			list[i].Name,
			list[i].Type,
			list[i].LanguageString(),
			emoji, list[i].Coverage, "%",
			list[i].ReportUrl)
	}
	if markdownStr != "" {
		markdownStr = fmt.Sprintf(
			`
## unit test report table 
| Service | Type | Languge | Total Project Coverage | Report Link |
| :--- | :--- | :--- | :--- | :--- |
%s`,
			markdownStr)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	_, _, err = client.Issues.CreateComment(context.Background(),
		strings.Split(repoPath, "/")[0],
		strings.Split(repoPath, "/")[1],
		number, &github.IssueComment{
			Body: &markdownStr,
		})
	return
}
