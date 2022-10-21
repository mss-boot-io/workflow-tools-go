/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/5/15 14:55
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/5/15 14:55
 */

package dep

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Report struct {
	Group              string `csv:"GROUP"`
	Package            string `csv:"PACKAGE"`
	Class              string `csv:"CLASS"`
	InstructionMissed  int    `csv:"INSTRUCTION_MISSED"`
	InstructionCovered int    `csv:"INSTRUCTION_COVERED"`
	BranchMissed       int    `csv:"BRANCH_MISSED"`
	BranchCovered      int    `csv:"BRANCH_COVERED"`
	LineMissed         int    `csv:"LINE_MISSED"`
	LineCovered        int    `csv:"LINE_COVERED"`
	ComplexityMissed   int    `csv:"COMPLEXITY_MISSED"`
	ComplexityCovered  int    `csv:"COMPLEXITY_COVERED"`
	MethodMissed       int    `csv:"METHOD_MISSED"`
	MethodCovered      int    `csv:"METHOD_COVERED"`
}

// GetAllServices 获取所有项目
func GetAllServices(workspace, filename, projectNameMatch string) (map[string][]string, error) {
	var services = make(map[string][]string)
	err := filepath.WalkDir(workspace, func(path string, d fs.DirEntry, err error) error {
		if filepath.Base(path) == filename {
			//找到目标文件
			projectName, _, err := GetGradleInfoFromPath(path, projectNameMatch, "")
			if err != nil {
				return err
			}
			services[projectName] = strings.Split(path[len(workspace)+1:len(path)-len(filename)-1], string(os.PathSeparator))
		}
		return nil
	})
	return services, err
}

// GetGradleInfoFromPath 获取项目gradle文件信息
func GetGradleInfoFromPath(path,
	projectNameMatch,
	dependenceMatch string) (projectName string, dependence []string, err error) {
	rb, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	r, err := regexp.Compile(projectNameMatch)
	if err != nil {
		return
	}
	//提取项目名称
	projectName = strings.Trim(strings.ReplaceAll(
		strings.ReplaceAll(string(r.Find(rb)), "'", ""),
		"rootProject.name = ",
		""), " ")
	if projectName == "" {
		dir, _ := filepath.Split(path)
		projectName = strings.Split(
			dir,
			string(os.PathSeparator))[len(strings.Split(dir, string(os.PathSeparator)))-2]
	}
	if dependenceMatch == "" {
		return
	}
	r, err = regexp.Compile(dependenceMatch)
	if err != nil {
		return
	}
	dependence = make([]string, 0)
	for _, dep := range r.FindAll(rb, -1) {
		dependence = append(dependence, strings.Trim(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(string(dep), "'", ""),
					"includeBuild ", ""),
				"../", ""),
			" "))
	}
	return
}
