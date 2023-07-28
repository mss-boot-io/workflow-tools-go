/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/10/8 15:57:10
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/10/8 15:57:10
 */

package cdk8s

import (
	"github.com/mss-boot-io/cd-template/pkg/config"
	"github.com/mss-boot-io/cd-template/stage"
)

// Generate the k8s yaml file
func Generate(configPath, stageStr, image string, servicePath []string) {
	config.NewConfig(&configPath)
	switch stageStr {
	case "dev", "test", "local", "alpha", "beta", "staging":
		config.Cfg.Hpa = false
		config.Cfg.Resources = nil
		config.Cfg.Replicas = config.Cfg.TestReplicas
	}
	config.Cfg.Image.Path = image
	stage.Synth(stageStr, servicePath...)
}
