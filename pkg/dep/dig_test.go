/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/4/7 16:29
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/4/7 16:29
 */

package dep

import (
	"reflect"
	"testing"

	"github.com/sanity-io/litter"
)

func TestNewDig(t *testing.T) {
	type args struct {
		workspace        string
		filename         string
		projectNameMatch string
		dependenceMatch  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test0",
			args{
				"../../testdata",
				"settings.gradle",
				"rootProject.name =\\s'([^']+)'",
				"includeBuild\\s'([^']+)'",
			},
			false,
		},
	}
	for _, tt := range tests {
		services, err := GetAllServices(tt.args.workspace, tt.args.filename, tt.args.projectNameMatch)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. GetAllServices() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDig(tt.args.workspace, tt.args.filename, services, tt.args.projectNameMatch, tt.args.dependenceMatch)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDig_GetChanged(t *testing.T) {
	services, err := GetAllServices("../../testdata",
		"settings.gradle",
		"rootProject.name =\\s'([^']+)'")
	if err != nil {
		t.Errorf("GetAllServices error = %v", err)
	}
	d, err := NewDig("../../testdata",
		"settings.gradle",
		services,
		"rootProject.name =\\s'([^']+)'",
		"includeBuild\\s'([^']+)'")
	if err != nil {
		t.Errorf("Init() error = %v", err)
	}

	type args struct {
		dirs []string
	}
	tests := []struct {
		name string
		d    *Dig
		args args
		want bool
	}{
		{
			"test0",
			d,
			args{
				[]string{"project0/lib0", "project2/lib2"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.GetChanged(tt.args.dirs)
			if !reflect.DeepEqual(got == nil, tt.want) {
				t.Errorf("GetChanged() = %v, want %v", got, tt.want)
			}
			litter.Dump(got)
		})
	}
}
