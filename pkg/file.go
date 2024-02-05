/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/4/7 16:19
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/4/7 16:19
 */

package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

// PathExist 判断目录或文件是否存在
func PathExist(addr string) bool {
	_, err := os.Stat(addr)
	if err != nil {
		return false
	}
	return true
}

// CreatePath 创建目录
func CreatePath(path string) error {
	if !PathExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteJsonFile 写入json文件
func WriteJsonFile(path string, data interface{}) (err error) {
	buffer := &bytes.Buffer{}
	switch data.(type) {
	case *bytes.Buffer:
		buffer = data.(*bytes.Buffer)
	case string:
		buffer.WriteString(data.(string))
	case []byte:
		buffer.Write(data.([]byte))
	default:
		err = json.NewEncoder(buffer).Encode(data)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(buffer.Bytes())
	return err
}

// ReadJsonFile 读取json文件
func ReadJsonFile(path string, data interface{}) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(data)
}

// CopyFile 复制文件
func CopyFile(sourceFile, destinationFile string) error {
	// open source file
	src, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer src.Close()

	//create destination file
	dst, err := os.Create(destinationFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	// copy file content
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

// ReadYamlFile 读取yaml文件
func ReadYamlFile(path string, data any) error {
	if !PathExist(path + ".yml") {
		if !PathExist(path + ".yaml") {
			return fmt.Errorf("%s not found", path)
		} else {
			path = path + ".yaml"
		}
	} else {
		path = path + ".yml"
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewDecoder(f).Decode(data)
}
