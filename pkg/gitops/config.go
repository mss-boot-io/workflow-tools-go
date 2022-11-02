/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/11/1 11:02:59
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/11/1 11:02:59
 */

package gitops

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config : config
type Config struct {
	Image   string `yaml:"image" json:"image"`
	Project string `yaml:"project" json:"project"`
	Deploy  Deploy `yaml:"deploy" json:"deploy"`
}

type Deploy struct {
	Stage map[string]StageDeploy `yaml:"stage" json:"stage"`
}

// StageDeploy : stage deploy
type StageDeploy struct {
	Cluster   string `yaml:"cluster" json:"cluster"`
	Namespace string `yaml:"namespace" json:"namespace"`
	AutoSync  bool   `yaml:"autoSync" json:"autoSync"`
}

// GetImage : get image
func (c *Config) GetImage(service string) string {
	if len(strings.Split(c.Image, "/")) > 1 {
		return c.Image
	}
	return fmt.Sprintf("%s/%s", c.Image, service)
}

// LoadFile : load file
func LoadFile(path string) (*Config, error) {
	config := &Config{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
