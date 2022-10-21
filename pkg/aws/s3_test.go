/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/4/8 14:36
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/4/8 14:36
 */

package aws

//import (
//	"github.com/sanity-io/litter"
//	"testing"
//
//	"github.com/mss-boot-io/workflow-tools/pkg/dep"
//)
//
//func TestGetObjectFromS3(t *testing.T) {
//	data := make([]*dep.Matrix, 0)
//	type args struct {
//		region string
//		bucket string
//		key    string
//		data   interface{}
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		{
//			"test0",
//			args{
//				region: "ap-northeast-1",
//				bucket: "github-runner-artifact",
//				key:    "testdata/test.json",
//				data:   &data,
//			},
//			false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := GetObjectFromS3(tt.args.region, tt.args.bucket, tt.args.key, tt.args.data); (err != nil) != tt.wantErr {
//				t.Errorf("GetObjectFromS3() error = %v, wantErr %v", err, tt.wantErr)
//			}
//			litter.Dump(data)
//		})
//	}
//}
