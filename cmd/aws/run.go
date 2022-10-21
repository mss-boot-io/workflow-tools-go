/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/4/13 11:40
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/4/13 11:40
 */

package aws

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/mss-boot-io/workflow-tools/pkg/aws"
)

var (
	outputType, bucket, region string
	StartCmd                   = &cobra.Command{
		Use:          "aws",
		Short:        "aws operator",
		Example:      "go-workflow-tools aws",
		SilenceUsage: true,
		PreRun: func(_ *cobra.Command, _ []string) {
			log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVar(&bucket,
		"bucket",
		os.Getenv("bucket"),
		"bucket")
	StartCmd.PersistentFlags().StringVar(&region,
		"region",
		os.Getenv("region"),
		"region")
	StartCmd.PersistentFlags().StringVar(&outputType,
		"output_type",
		os.Getenv("output_type"),
		"output type")
	if outputType == "" {
		outputType = "ec2"
	}
}

func run() error {
	switch outputType {
	case "ec2":
		return outputEC2List()
	case "rds":
		return outputRDSList()
	}
	return fmt.Errorf("output type %s not support", outputType)
}

func outputEC2List() error {
	records := [][]string{
		{"Region", "Name", "InstanceId", "InstanceType", "ClusterName", "State"},
	}
	regions := []string{"ap-northeast-1", "us-west-2"}
	for i := range regions {
		fmt.Printf("start to get %s ec2 list...\n", regions[i])
		rs, err := aws.GetEC2List(regions[i])
		if err != nil {
			return err
		}
		records = append(records, rs...)
		fmt.Printf("finish get %s ec2 list.\n", regions[i])
	}
	//统计ri
	records = append(records, make([][]string, 1)...)
	records = append(records, []string{"Region", "ReservedInstanceId", "InstanceType", "InstanceCount", "State"})

	for i := range regions {
		fmt.Printf("start to get %s ri list...\n", regions[i])
		ris, err := aws.GetEC2RIList(regions[i])
		if err != nil {
			return err
		}
		records = append(records, ris...)
		fmt.Printf("finish get %s ri listn.\n", regions[i])
	}

	buffer := &bytes.Buffer{}
	w := csv.NewWriter(buffer)
	defer w.Flush()
	err := w.WriteAll(records)
	if err != nil {
		return err
	}
	return aws.PutObjectToS3(region, bucket,
		fmt.Sprintf("aws/ec2-%s.csv",
			time.Now().Format("2006-01-02")),
		buffer, "text/csv")
}

func outputRDSList() error {
	records := [][]string{
		{"Region", "Name", "Engine", "InstanceType", "Status"},
	}
	regions := []string{"ap-northeast-1", "us-west-2"}
	// 统计RDS
	for i := range regions {
		fmt.Printf("start to get %s rds list...\n", regions[i])
		rs, err := aws.GetRDSList(regions[i])
		if err != nil {
			return err
		}
		records = append(records, rs...)
		fmt.Printf("finish get %s rds list.\n", regions[i])
	}
	// 统计RDS RI
	records = append(records, make([][]string, 1)...)
	records = append(records, []string{"Region", "ReservedInstanceId", "InstanceType", "InstanceCount", "State"})

	for i := range regions {
		fmt.Printf("start to get %s rds ri list...\n", regions[i])
		ris, err := aws.GetRDSReservedList(regions[i])
		if err != nil {
			return err
		}
		records = append(records, ris...)
		fmt.Printf("finish get %s rds ri list...\n", regions[i])
	}

	buffer := &bytes.Buffer{}
	w := csv.NewWriter(buffer)
	defer w.Flush()
	err := w.WriteAll(records)
	if err != nil {
		return err
	}
	return aws.PutObjectToS3(region, bucket,
		fmt.Sprintf("aws/rds-%s.csv",
			time.Now().Format("2006-01-02")),
		buffer, "text/csv")
}
