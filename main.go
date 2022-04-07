/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/4/6 11:18
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/4/6 11:18
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mss-boot-io/multiple-work/work"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	leafs := make([]work.Leaf, 0)
	err := json.Unmarshal([]byte(os.Getenv("leaf")), &leafs)
	if err != nil {
		log.Println(err)
		return
	}

	for i := range leafs {
		fmt.Printf("######################## %s ########################\n", leafs[i].Name)
		leafs[i].Err = leafs[i].Run(os.Getenv("cmd"))
		fmt.Print("###   ")
		if leafs[i].Err != nil {
			fmt.Println("Failed")
		} else {
			fmt.Println("Successful")
		}
		fmt.Printf("######################## %s ########################\n", leafs[i].Name)
	}

	fmt.Println()
	fmt.Println()
	fmt.Println()

	fmt.Printf("######################## %s ########################\n", "All Service")
	var failed bool
	for i := range leafs {
		fmt.Printf("###   %s: ", leafs[i].Name)
		if leafs[i].Err != nil {
			failed = true
			fmt.Println("Failed")
			continue
		}
		fmt.Println("Successful")
	}
	fmt.Printf("######################## %s ########################\n", "All Service")
	if failed {
		os.Exit(-1)
	}
}
