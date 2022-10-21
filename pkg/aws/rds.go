package aws

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// GetRDSList 获取RDS实例列表
func GetRDSList(region string) ([][]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	client := rds.NewFromConfig(cfg)
	input := &rds.DescribeDBInstancesInput{}

	result, err := GetRDSInstances(context.TODO(), client, input)

	records := make([][]string, 0)
	for _, instance := range result.DBInstances {
		records = append(records, []string{
			region,
			*instance.DBInstanceIdentifier,
			*instance.Engine,
			*instance.DBInstanceClass,
			*instance.DBInstanceStatus,
		})
	}
	return records, nil
}

// GetRDSReservedList 获取RDS RI列表
func GetRDSReservedList(region string) ([][]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	client := rds.NewFromConfig(cfg)
	input := &rds.DescribeReservedDBInstancesInput{}

	result, err := GetRDSReservedInstances(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error retrieving information about your rds reserved instances:", err)
		return nil, err
	}
	records := make([][]string, 0)
	for _, ri := range result.ReservedDBInstances {
		records = append(records, []string{
			region,
			*ri.ReservedDBInstanceId,
			*ri.DBInstanceClass,
			strconv.Itoa(int(ri.DBInstanceCount)),
			*ri.State,
		})
	}
	return records, nil
}

type RDSDescribeInstancesAPI interface {
	DescribeDBInstances(ctx context.Context,
		params *rds.DescribeDBInstancesInput,
		optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
}

func GetRDSInstances(c context.Context, api RDSDescribeInstancesAPI,
	input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
	return api.DescribeDBInstances(c, input)
}

type RDSDescribeReservedInstancesAPI interface {
	DescribeReservedDBInstances(ctx context.Context,
		params *rds.DescribeReservedDBInstancesInput,
		optFns ...func(*rds.Options)) (*rds.DescribeReservedDBInstancesOutput, error)
}

func GetRDSReservedInstances(c context.Context, api RDSDescribeReservedInstancesAPI,
	input *rds.DescribeReservedDBInstancesInput) (*rds.DescribeReservedDBInstancesOutput, error) {
	return api.DescribeReservedDBInstances(c, input)
}
