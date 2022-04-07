/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/4/6 11:19
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/4/6 11:19
 */

package work

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Leaf struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Err  error  `json:"-"`
}

func (e *Leaf) Run(cs string) error {
	cs = fmt.Sprintf("cd %s && %s", e.Name, cs)
	cmd := exec.Command("/bin/bash", "-c", cs)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
