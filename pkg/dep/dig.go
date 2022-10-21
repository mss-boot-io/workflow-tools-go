/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/4/7 16:11
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/4/7 16:11
 */

package dep

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mss-boot-io/workflow-tools/pkg"
)

// Dig dependency dig
type Dig struct {
	data             map[string][]string
	workspace        string
	filename         string
	projectNameMatch string
	dependenceMatch  string
	servicePath      map[string][]string
}

// NewDig new dig
func NewDig(workspace, filename string, servicePath map[string][]string, projectNameMatch, dependenceMatch string) (*Dig, error) {
	d := &Dig{
		data:             make(map[string][]string),
		servicePath:      servicePath,
		workspace:        workspace,
		filename:         filename,
		projectNameMatch: projectNameMatch,
		dependenceMatch:  dependenceMatch,
	}
	err := d.init()
	return d, err
}

// init dig
func (d *Dig) init() error {
	if !pkg.PathExist(d.workspace) {
		return errors.New("workspace not exist")
	}
	if d == nil {
		return errors.New("dig is nil")
	}
	return filepath.WalkDir(d.workspace,
		d.walkDirHandler(d.filename, d.projectNameMatch, d.dependenceMatch))
}

func (d *Dig) get(key string) ([]string, bool) {
	if d == nil {
		return nil, false
	}
	v, ok := d.data[d.getServiceByPath(key)]
	return v, ok
}

func (d *Dig) set(key string, value ...string) {
	if d.data == nil {
		d.data = make(map[string][]string)
	}
	v, ok := d.get(key)
	if ok {
		d.data[key] = append(v, value...)
		return
	}
	d.data[key] = value
}

func (d *Dig) getMatrix(data map[string]*Matrix) []*Matrix {
	matrix := make([]*Matrix, 0, len(d.data))
	for _, m := range data {
		matrix = append(matrix, m)
	}
	return matrix
}

// GetChanged get changed
func (d *Dig) GetChanged(dirs []string) []*Matrix {
	matrix := make(map[string]*Matrix)
	for i := range dirs {
		// get library
		if _, ok := d.get(dirs[i]); ok {
			matrix[dirs[i]] = &Matrix{
				Name:        d.getServiceByPath(dirs[i]),
				Type:        Library,
				ProjectPath: strings.Split(dirs[i], "/"),
			}
		}
	}
	// get leaf
	for _, m := range d.getLeaf(dirs) {
		matrix[m.Name] = m
	}
	return d.getMatrix(matrix)
}

// getLeaf get leaf
func (d *Dig) getLeaf(dirs []string) map[string]*Matrix {
	matrix := make(map[string]*Matrix)
	for i := range dirs {
		if children, ok := d.get(dirs[i]); !ok {
			if strings.Index(strings.ToLower(dirs[i]), Lambda.String()) == -1 &&
				strings.Index(strings.ToLower(dirs[i]), Airflow.String()) == -1 {
				// 不是lambda和airflow
				//判断是否存在Dockerfile文件
				if pkg.PathExist(filepath.Join(d.workspace, dirs[i], "Dockerfile")) {
					if os.Getenv("pr_num") != "" {
						matrix[dirs[i]] = &Matrix{
							Name:      d.getServiceByPath(dirs[i]),
							Type:      Service,
							ReportUrl: os.Getenv("cloudfront_url") + "/" + os.Getenv("github_repository") + "/" + "pr" + os.Getenv("pr_num") + "/" + dirs[i] + "/jacoco-report/index.html",
						}
					} else {
						matrix[dirs[i]] = &Matrix{
							Name: d.getServiceByPath(dirs[i]),
							Type: Service,
						}
					}
				}
				continue
			}
			if strings.Index(strings.ToLower(dirs[i]), Lambda.String()) != -1 {
				// 是lambda
				matrix[dirs[i]] = &Matrix{
					Name: d.getServiceByPath(dirs[i]),
					Type: Service,
				}
				continue
			}
			if strings.Index(strings.ToLower(dirs[i]), Airflow.String()) != -1 {
				// 是airflow
				matrix[dirs[i]] = &Matrix{
					Name: d.getServiceByPath(dirs[i]),
					Type: Service,
				}
				continue
			}
			return matrix
		} else {
			for _, m := range d.getLeaf(children) {
				matrix[m.Name] = m
			}
		}
	}
	return matrix
}

// walkDirHandler walk dir
func (d *Dig) walkDirHandler(filename, projectNameMatch, dependenceMatch string) fs.WalkDirFunc {
	return func(path string, dir fs.DirEntry, err error) error {
		if dir.IsDir() {
			return nil
		}
		if filepath.Base(path) == filename {
			//获取到目标文件
			projectName, dependence, err := GetGradleInfoFromPath(path, projectNameMatch, dependenceMatch)
			if err != nil {
				return err
			}
			for i := range dependence {
				d.set(dependence[i], projectName)
			}
		}
		return nil
	}
}

func (d *Dig) getServiceByPath(path string) string {
	for service := range d.servicePath {
		if strings.Join(d.servicePath[service], "/") == path {
			return service
		}
	}
	return ""
}
