package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
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

func GetEC2Client() *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error: Loading the config")
		log.Fatal(err)
	}
	return ec2.NewFromConfig(cfg)
}

func startTask() {
	svc := GetECSClient()
	ec2Svc := GetEC2Client()
	cluster := aws.String("streamyard")
	taskDef := aws.String("streamer")
	input := &ecs.RunTaskInput{
		Cluster:        cluster,
		TaskDefinition: taskDef,
		LaunchType:     types.LaunchTypeFargate,
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        []string{"your-subnet-id"},
				AssignPublicIp: types.AssignPublicIpEnabled,
				SecurityGroups: []string{"your-security-group-id"},
			},
		},
	}
	result, err := svc.RunTask(context.TODO(), input)
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}

	if len(result.Tasks) == 0 {
		log.Println("No Task started")
		return
	}

	task := result.Tasks[0]
	fmt.Printf("Started task: %s\n", *task.TaskArn)

	describeTaskInput := &ecs.DescribeTasksInput{
		Cluster: cluster,
		Tasks:   []string{*task.TaskArn},
	}

	describeTasksResult, err := svc.DescribeTasks(context.TODO(), describeTaskInput)
	if err != nil {
		log.Fatal("error describing the task: ", err)
		return
	}

	if len(describeTasksResult.Tasks) == 0 || len(describeTasksResult.Tasks[0].Attachments) == 0 {
		log.Println("No task descriptions found")
		return
	}

	var eniID string
	for _, detail := range describeTasksResult.Tasks[0].Attachments[0].Details {
		if aws.ToString(detail.Name) == "networkInterfaceId" {
			eniID = aws.ToString(detail.Value)
			break
		}
	}

	if eniID == "" {
		log.Println("No ENI ID found for the task")
		return
	}

	// Describe the ENI to get the public IP
	describeNetworkInterfacesInput := &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []string{eniID},
	}

	describeNetworkInterfacesResult, err := ec2Svc.DescribeNetworkInterfaces(context.TODO(), describeNetworkInterfacesInput)
	if err != nil {
		log.Fatal("Error describing network interface: ", err)
		return
	}

	if len(describeNetworkInterfacesResult.NetworkInterfaces) == 0 {
		log.Println("No network interfaces found")
		return
	}

	publicIP := describeNetworkInterfacesResult.NetworkInterfaces[0].Association.PublicIp
	fmt.Printf("Task public IP: %s\n", *publicIP)
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
