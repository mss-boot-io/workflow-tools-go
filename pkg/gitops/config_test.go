/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/11/2 10:54:25
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/11/2 10:54:25
 */

package gitops

import (
	"github.com/sanity-io/litter"
	"path/filepath"
	"testing"
)

func TestLoadFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test load file",
			args: args{
				path: filepath.Join("testdata", "config.yml"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			litter.Dump(got)
		})
	}
}
