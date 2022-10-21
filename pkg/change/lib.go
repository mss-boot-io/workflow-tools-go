/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/8/30 16:18:01
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/8/30 16:18:01
 */

package change

import (
	"fmt"
)

// GetFilename get filename for change files list
func GetFilename(repo, mark string) string {
	return fmt.Sprintf("%s/%s/artifact/workflow/changed.json", repo, mark)
}
