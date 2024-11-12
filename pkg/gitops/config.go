/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/11/1 11:02:59
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/11/1 11:02:59
 */

package gitops

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Config : config
type Config struct {
	Project            string `yaml:"project" json:"project"`
	LanguageEnvType    string `yaml:"languageEnvType" json:"languageEnvType"`
	LanguageEnvVersion string `yaml:"languageEnvVersion" json:"languageEnvVersion"`
	LanguageEnvCache   string `yaml:"languageEnvCache" json:"languageEnvCache"`
	ArmImageNeeds      bool   `yaml:"armImageNeeds" json:"armImageNeeds"`
	Deploy             Deploy `yaml:"deploy" json:"deploy"`
	Build              Build  `yaml:"build" json:"build"`
}

type Deploy struct {
	Image string                 `yaml:"image" json:"image"`
	Stage map[string]StageDeploy `yaml:"stage" json:"stage"`
}

type Build struct {
	SkipCache bool `yaml:"skipCache" json:"skipCache"`
}

type StageDeploy map[string]any

func (s StageDeploy) GetKey(key string) any {
	v, ok := s[key]
	if !ok {
		return ""
	}
	return v
}

// ParseTemplate : parse template
func (s StageDeploy) ParseTemplate(tmp string) (string, error) {
	var err error
	t := template.New(tmp)
	t, err = t.Parse(tmp)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = t.Execute(&buffer, s)
	return buffer.String(), err
}

// StageDeploy : stage deploy
//type StageDeploy struct {
//	Cluster   string `yaml:"cluster" json:"cluster"`
//	Namespace string `yaml:"namespace" json:"namespace"`
//	AutoSync  bool   `yaml:"autoSync" json:"autoSync"`
//}

// GetImage : get image
func (c *Config) GetImage(service string) string {
	if c.Deploy.Image == "" {
		return ""
	}
	if (!strings.Contains(c.Deploy.Image, "public.ecr.aws") &&
		len(strings.Split(c.Deploy.Image, "/")) > 1) ||
		(strings.Contains(c.Deploy.Image, "public.ecr.aws") &&
			len(strings.Split(c.Deploy.Image, "/")) > 2) {
		return c.Deploy.Image
	}
	return fmt.Sprintf("%s/%s", c.Deploy.Image, service)
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
