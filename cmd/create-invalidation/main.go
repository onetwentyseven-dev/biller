package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	if len(os.Args) != 2 {
		logger.Fatal("len of args should be 1")
	}

	distributionID := os.Args[1]

	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.WithError(err).Fatal("failed to install aws config")
	}

	cloudfrontClient := cloudfront.NewFromConfig(awsCfg)

	getDistributionOutput, err := cloudfrontClient.GetDistribution(ctx, &cloudfront.GetDistributionInput{
		Id: aws.String(distributionID),
	})
	if err != nil {
		logger.WithError(err).Fatal("failed to fetch distribution")
	}

	distribution := getDistributionOutput.Distribution

	logger.Info("Validated Distribution")

	createInvalidationOutput, err := cloudfrontClient.CreateInvalidation(ctx, &cloudfront.CreateInvalidationInput{
		DistributionId: distribution.Id,
		InvalidationBatch: &types.InvalidationBatch{
			CallerReference: aws.String(fmt.Sprintf("OTS Frontend Deployment Script: %s", time.Now())),
			Paths: &types.Paths{
				Items:    []string{"/*"},
				Quantity: aws.Int32(1),
			},
		},
	})
	if err != nil {
		logger.WithError(err).Fatal("failed to create invalidation")
	}

	invalidation := createInvalidationOutput.Invalidation

	invalidationData, _ := json.Marshal(invalidation)
	fmt.Println("invalidation :: ", string(invalidationData))

	logger.WithField("InvalidationID", aws.ToString(invalidation.Id)).Info("Invalidation Created Successfully")

	for {
		getInvalidation, err := cloudfrontClient.GetInvalidation(ctx, &cloudfront.GetInvalidationInput{
			DistributionId: distribution.Id,
			Id:             invalidation.Id,
		})
		if err != nil {
			logger.WithError(err).Fatal("failed to get invalidation status")
		}

		invalidation := getInvalidation.Invalidation

		status := aws.ToString(invalidation.Status)
		entry := logger.WithField("status", status)
		if status == "Completed" {
			entry.Info("invalidation done")
			break
		}

		entry.Info("invalidation in progress....sleeping for 10 seconds")

		time.Sleep(time.Second * 10)
	}

}
