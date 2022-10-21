/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/4/13 11:47
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/4/13 11:47
 */

package aws

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// GetEC2RIList 获取EC2 RI列表
func GetEC2RIList(region string) ([][]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	client := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeReservedInstancesInput{}

	result, err := GetReservedInstances(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error retrieving information about your Amazon EC2 reserved instances:", err)
		return nil, err
	}
	records := make([][]string, 0)
	for _, ri := range result.ReservedInstances {
		records = append(records, []string{
			region,
			*ri.ReservedInstancesId,
			string(ri.InstanceType),
			strconv.Itoa(int(*ri.InstanceCount)),
			string(ri.State),
		})
	}
	return records, nil
}

// GetEC2List returns an EC2 list for region
func GetEC2List(region string) ([][]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	client := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeInstancesInput{}

	result, err := GetInstances(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error retrieving information about your Amazon EC2 instances:", err)
		return nil, err
	}

	records := make([][]string, 0)

	for _, r := range result.Reservations {
		for i := range r.Instances {
			clusterName := ""
			name := ""
			for j := range r.Instances[i].Tags {
				if *(r.Instances[i].Tags[j].Key) == "Name" {
					name = *(r.Instances[i].Tags[j].Value)
					continue
				}
				if *(r.Instances[i].Tags[j].Key) == "eks:cluster-name" {
					//cluster-name
					clusterName = *(r.Instances[i].Tags[j].Value)
					continue
				}
			}
			records = append(records, []string{
				region,
				name,
				*r.Instances[i].InstanceId,
				string(r.Instances[i].InstanceType),
				clusterName,
				string(r.Instances[i].State.Name),
			})
		}
	}
	return records, nil
}

// EC2DescribeInstancesAPI defines the interface for the DescribeInstances function.
// We use this interface to test the function using a mocked service.
type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

type EC2ReservedInstancesAPI interface {
	DescribeReservedInstances(ctx context.Context,
		params *ec2.DescribeReservedInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeReservedInstancesOutput, error)
}

// GetInstances retrieves information about your Amazon Elastic Compute Cloud (Amazon EC2) instances.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a DescribeInstancesOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to DescribeInstances.
func GetInstances(c context.Context, api EC2DescribeInstancesAPI,
	input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}

func GetReservedInstances(c context.Context, api EC2ReservedInstancesAPI,
	input *ec2.DescribeReservedInstancesInput) (*ec2.DescribeReservedInstancesOutput, error) {
	return api.DescribeReservedInstances(c, input)
}
