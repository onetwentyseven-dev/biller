package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	asTypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/sirupsen/logrus"
)

func main() {

	targetASG := "ots-bastion"

	logger := logrus.New()

	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.WithError(err).Fatal("failed to initiali aws config")
	}

	scalingClient := autoscaling.NewFromConfig(awsCfg)
	ec2Client := ec2.NewFromConfig(awsCfg)

	groups, err := scalingClient.DescribeAutoScalingGroups(ctx, &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{
			targetASG,
		},
	})
	if err != nil {
		logger.WithError(err).Fatal("failed to fetch autoscaling groups")
	}

	var group *asTypes.AutoScalingGroup
	for _, asGroup := range groups.AutoScalingGroups {
		if aws.ToString(asGroup.AutoScalingGroupName) == targetASG {
			group = &asGroup
			break
		}
	}

	if group == nil {
		logger.Fatal("Target ASG was not found is results")
	}

	if len(group.Instances) == 0 {
		fmt.Printf("ASG does not have any instances. Update Desired Capcity to spin one up? [y/n]: ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if text == "n" {
			fmt.Println("Receive n response, exiting script")
			os.Exit(0)
		}

		if text != "y" {
			fmt.Println("Invalid response received, exiting script")
			os.Exit(1)
		}

		_, err := scalingClient.SetDesiredCapacity(ctx, &autoscaling.SetDesiredCapacityInput{
			AutoScalingGroupName: aws.String(targetASG),
			DesiredCapacity:      aws.Int32(1),
		})
		if err != nil {
			fmt.Printf("Failed to update desired capacity: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("Successfully update auto scaling group, launching waiter")
		time.Sleep(time.Second * 10)

		groups, err = scalingClient.DescribeAutoScalingGroups(ctx, &autoscaling.DescribeAutoScalingGroupsInput{
			AutoScalingGroupNames: []string{targetASG},
		})
		if err != nil {
			fmt.Printf("Failed to update desired capacity: %s\n", err)
			os.Exit(1)
		}

		for _, asGroup := range groups.AutoScalingGroups {
			if aws.ToString(asGroup.AutoScalingGroupName) == targetASG {
				group = &asGroup
				break
			}
		}

	}

	if len(group.Instances) > 1 {
		fmt.Println("NUMBER OF INSTANCES IS GREATER THAN 1. CHECK THE CONSOLE FOR MORE DETAILS")
		os.Exit(1)
	}

	fmt.Println("Using Ec2 Client to check status of instance and grab public ip address")

	instanceID := aws.ToString(group.Instances[0].InstanceId)

	waiter := ec2.NewInstanceRunningWaiter(ec2Client)
	output, err := waiter.WaitForOutput(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}, time.Minute*5)
	if err != nil {
		fmt.Printf("Failed to update desired capacity: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%d Reservations Returned\n", len(output.Reservations))

	for _, res := range output.Reservations {
		for _, instance := range res.Instances {
			for _, netInterfaces := range instance.NetworkInterfaces {
				fmt.Printf("%s - %s\n", aws.ToString(instance.InstanceId), aws.ToString(netInterfaces.Association.PublicIp))
			}
		}
	}

}
