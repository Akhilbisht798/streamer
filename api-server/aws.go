package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func GetECSClient() *ecs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error: Loading the config")
		log.Fatal(err)
	}
	svc := ecs.NewFromConfig(cfg)
	return svc
}

func startTask() (string, string) {
	svc := GetECSClient()

	// Configuration for the task
	taskDefination := aws.String("")
	cluster := aws.String("")
	networkConfigure := types.NetworkConfiguration{
		AwsvpcConfiguration: &types.AwsVpcConfiguration{
			Subnets:        []string{*aws.String(""), *aws.String(""), *aws.String("")},
			SecurityGroups: []string{*aws.String("")},
			AssignPublicIp: types.AssignPublicIpEnabled,
		},
	}
	containerOverride := []types.ContainerOverride{
		{
			Name: aws.String(""),
			Environment: []types.KeyValuePair{
				{
					Name:  aws.String("LINK"),
					Value: aws.String(""),
				},
			},
		},
	}

	input := &ecs.RunTaskInput{
		TaskDefinition:       taskDefination,
		Cluster:              cluster,
		LaunchType:           types.LaunchTypeFargate,
		NetworkConfiguration: &networkConfigure,
		Overrides: &types.TaskOverride{
			ContainerOverrides: containerOverride,
		},
	}

	// Starting the task.
	result, err := svc.RunTask(context.TODO(), input)
	if err != nil {
		log.Fatal(err)
	}

	// Waiting for task to start.
	// describeInputTask := &ecs.DescribeTasksInput{
	// 	Tasks:   []string{*result.Tasks[0].TaskArn},
	// 	Cluster: cluster,
	// }
	// for {
	// 	res, err := svc.DescribeTasks(context.TODO(), describeInputTask)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 		break
	// 	}
	// 	fmt.Println("Pending State: ", res.Tasks[0].TaskArn)
	// 	if *res.Tasks[0].TaskArn == *result.Tasks[0].TaskArn && *res.Tasks[0].LastStatus == "RUNNING" {
	// 		fmt.Println("Task Status Running")
	// 		break
	// 	}
	// 	time.Sleep(2 * time.Second)
	// }

	// getting the public ip addr.
	var publicIp string
	describeInputTask := &ecs.DescribeTasksInput{
		Tasks:   []string{*result.Tasks[0].TaskArn},
		Cluster: cluster,
	}
	res, err := svc.DescribeTasks(context.TODO(), describeInputTask)
	if err != nil {
		log.Fatal(err)
	}
	for _, task := range res.Tasks {
		for _, attachment := range task.Attachments {
			if *attachment.Type == "ElasticNetworkInterface" {
				for _, detail := range attachment.Details {
					if *detail.Name == "networkInterfaceId" {
						eniID := *detail.Value
						fmt.Println("ENI ID:", eniID)
					}
				}
			}
		}
	}

	return *result.Tasks[0].TaskArn, publicIp
}

func stopTask(taskArn string) {
	svc := GetECSClient()
	input := &ecs.StopTaskInput{
		Task:    aws.String(taskArn),
		Cluster: aws.String(""),
	}

	resp, err := svc.StopTask(context.TODO(), input)
	if err != nil {
		log.Fatal("Error Stopping task: ", err)
		return
	}

	log.Println("Task Stopped Successfully: ", resp)
}
