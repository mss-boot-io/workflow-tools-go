/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/5/19 11:04
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/5/19 11:04
 */

package dep

import "testing"

func TestOutputReportTableToPR(t *testing.T) {
	type args struct {
		repoPath string
		number   int
		list     []Matrix
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				repoPath: "mss-boot-io/workflow-tools-go",
				number:   1,
				list: []Matrix{
					{
						Name:        "service0",
						Type:        Service,
						Language:    []Language{"java", "go"},
						ReportUrl:   "https://github.com",
						ProjectPath: []string{"testdata", "project0", "service0"},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := OutputReportTableToPR(tt.args.repoPath, tt.args.number, tt.args.list); (err != nil) != tt.wantErr {
				t.Errorf("OutputReportTableToPR() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
