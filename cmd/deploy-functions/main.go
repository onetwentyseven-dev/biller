package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/sirupsen/logrus"
)

func main() {
	const directory = "./dist"

	logger := logrus.New()

	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.WithError(err).Fatal("failed to initiali aws config")
	}

	lambdaClient := lambda.NewFromConfig(awsCfg)

	files, err := os.ReadDir(directory)
	if err != nil {
		logger.WithError(err).Fatal("failed to list directory")
		return
	}

	var wg = new(sync.WaitGroup)

	for _, f := range files {
		wg.Add(1)
		if !strings.HasSuffix(f.Name(), ".zip") {
			continue
		}
		go func(dirEnt fs.DirEntry, wg *sync.WaitGroup) {
			defer wg.Done()
			entry := logger.WithField("file", dirEnt.Name())
			file, err := os.Open(filepath.Join(directory, dirEnt.Name()))
			if err != nil {
				entry.WithError(err).Error("failed to open zip")
				return
			}

			data, err := io.ReadAll(file)
			if err != nil {
				entry.WithError(err).Error("failed to read zip")
				return
			}

			handler := strings.TrimSuffix(dirEnt.Name(), ".zip")

			output, err := lambdaClient.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
				FunctionName: aws.String(fmt.Sprintf("%s_handler", handler)),
				ZipFile:      data,
			})
			if err != nil {
				entry.WithError(err).Error("failed to update function")
				return
			}

			entry.WithFields(logrus.Fields{
				"last_modified": aws.ToString(output.LastModified),
			}).Info("function updated successfully")

		}(f, wg)
	}

	logger.Info("updating functions")
	wg.Wait()

	logger.Info("done")

}
