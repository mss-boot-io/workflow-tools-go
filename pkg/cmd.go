/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/4/22 18:22
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/4/22 18:22
 */

package pkg

import (
	"os"
	"os/exec"
)

// Cmd 执行命令
func Cmd(cs string) error {
	cmd := exec.Command("/bin/bash", "-c", cs)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
