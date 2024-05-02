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

func startTask() {
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
