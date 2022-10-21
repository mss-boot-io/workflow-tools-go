/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/8/29 10:45:10
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/8/29 10:45:10
 */

package change

type Changer interface {
	// SetAuth 设置登录验证
	SetAuth(auth interface{})
	// SetRepoURL 设置代码仓库地址
	SetRepoURL(string) error
	// ChangeFiles 获取变化的文件列表
	ChangeFiles(string) (*Files, error)
}

type Files struct {
	Added    []string `json:"added"`
	Modified []string `json:"modified"`
	Deleted  []string `json:"deleted"`
	Renamed  []string `json:"renamed"`
}
