/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/11/1 11:02:59
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/11/1 11:02:59
 */

package gitops

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config : config
type Config struct {
	Project string                 `yaml:"project" json:"project"`
	Stage   map[string]StageDeploy `yaml:"deploy" json:"deploy"`
}

// StageDeploy : stage deploy
type StageDeploy struct {
	Cluster   string `yaml:"cluster" json:"cluster"`
	Namespace string `yaml:"namespace" json:"namespace"`
	AutoSync  bool   `yaml:"autoSync" json:"autoSync"`
}

// LoadFile : load file
func LoadFile(path string) (*Config, error) {
	config := struct {
		Deploy  Config `yaml:"deploy" json:"deploy"`
		Project string `yaml:"project" json:"project"`
	}{}
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}
	config.Deploy.Project = config.Project
	return &config.Deploy, nil
}
