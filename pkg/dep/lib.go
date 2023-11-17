/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/8/30 17:01:29
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/8/30 17:01:29
 */

package dep

import (
	"fmt"
	"path/filepath"
	"strings"
)

// GetFilename get filename for dep service list
func GetFilename(repo, mark, provider string) string {
	switch strings.ToLower(provider) {
	case "s3", "minio":
		return fmt.Sprintf("%s/%s/artifact/workflow/service.json", repo, mark)
	default:
		return filepath.Join("artifact", "workflow", "service.json")
	}
}
