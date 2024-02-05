package dep

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/mss-boot-io/workflow-tools/pkg"
	"gopkg.in/yaml.v3"
)

/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2024/2/5 12:01:46
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2024/2/5 12:01:46
 */

type GroupRule struct {
	Groups []*RuleGroupDesc `json:"groups"`
}

// RuleGroupDesc is a proto representation of a mimir rule group.
type RuleGroupDesc struct {
	Name      string        `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace string        `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Interval  time.Duration `protobuf:"bytes,3,opt,name=interval,proto3,stdduration" json:"interval"`
	Rules     []*RuleDesc   `protobuf:"bytes,4,rep,name=rules,proto3" json:"rules,omitempty"`
	User      string        `protobuf:"bytes,6,opt,name=user,proto3" json:"user,omitempty"`
	// The options field can be used to extend Mimir Ruler functionality without
	// having to repeatedly redefine the proto description. It can also be leveraged
	// to create custom `ManagerOpts` based on rule configs which can then be passed
	// to the Prometheus Manager.
	Options                       []*types.Any  `protobuf:"bytes,9,rep,name=options,proto3" json:"options,omitempty"`
	SourceTenants                 []string      `protobuf:"bytes,10,rep,name=sourceTenants,proto3" json:"sourceTenants,omitempty"`
	EvaluationDelay               time.Duration `protobuf:"bytes,11,opt,name=evaluationDelay,proto3,stdduration" json:"evaluationDelay"`
	AlignEvaluationTimeOnInterval bool          `protobuf:"varint,12,opt,name=align_evaluation_time_on_interval,json=alignEvaluationTimeOnInterval,proto3" json:"align_evaluation_time_on_interval,omitempty"`
}

// RuleDesc is a proto representation of a Prometheus Rule
type RuleDesc struct {
	Expr          string            `protobuf:"bytes,1,opt,name=expr,proto3" json:"expr,omitempty"`
	Record        string            `protobuf:"bytes,2,opt,name=record,proto3" json:"record,omitempty"`
	Alert         string            `protobuf:"bytes,3,opt,name=alert,proto3" json:"alert,omitempty"`
	For           time.Duration     `protobuf:"bytes,4,opt,name=for,proto3,stdduration" json:"for"`
	KeepFiringFor time.Duration     `protobuf:"bytes,13,opt,name=keep_firing_for,json=keepFiringFor,proto3,stdduration" json:"keep_firing_for"`
	Labels        map[string]string `protobuf:"bytes,5,rep,name=labels,proto3,customtype=github.com/grafana/mimir/pkg/mimirpb.LabelAdapter" json:"labels"`
	Annotations   map[string]string `protobuf:"bytes,6,rep,name=annotations,proto3,customtype=github.com/grafana/mimir/pkg/mimirpb.LabelAdapter" json:"annotations"`
}

func (e *Matrix) Alert(workspace, filename, prometheusPrefix, prometheusUsername, prometheusPassword string) error {
	if len(e.ProjectPath) == 0 {
		e.ProjectPath = []string{e.Name}
	}
	if workspace != "" {
		filename = filepath.Join(workspace, filepath.Join(e.ProjectPath...), filename)
	} else {
		filename = filepath.Join(filepath.Join(e.ProjectPath...), filename)
	}
	rules := &GroupRule{}
	err := pkg.ReadYamlFile(filename, rules)
	if err != nil {
		return err
	}
	group := "anonymous"
	if e.Project != "" {
		group = e.Project
	}
	url := fmt.Sprintf("%s/config/v1/rules/%s", prometheusPrefix, group)
	for i := range rules.Groups {
		r, _ := yaml.Marshal(rules.Groups[i])
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(r))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "text/plain")
		req.SetBasicAuth(prometheusUsername, prometheusPassword)
		client := &http.Client{}
		_, err = client.Do(req)
		if err != nil {
			return err
		}
	}
	return nil
}
