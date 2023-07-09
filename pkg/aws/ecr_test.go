package aws

import "testing"

/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2023/7/9 15:36:03
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2023/7/9 15:36:03
 */

func TestCreatePrivateRepoIfNotExist(t *testing.T) {
	type args struct {
		region string
		image  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				region: "us-west-2",
				image:  "316274061697.dkr.ecr.us-west-2.amazonaws.com/cluster-manager",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreatePrivateRepoIfNotExist(tt.args.region, tt.args.image); (err != nil) != tt.wantErr {
				t.Errorf("CreatePrivateRepoIfNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
