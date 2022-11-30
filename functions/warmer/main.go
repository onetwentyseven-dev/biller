package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/onetwentyseven-dev/apigw"
	"github.com/onetwentyseven-dev/biller/internal/mysql"
	"github.com/sirupsen/logrus"
)

type handler struct {
	logger *logrus.Logger
}

func main() {

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	awsCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize aws config")
	}

	loadConfig(awsCfg)

	api := apigw.New(logger)

	h := &handler{
		logger: logger,
	}

	api.AddHandler(http.MethodGet, "/warmer", h.handleGetWarmer)

	lambda.Start(api.HandleRoutes)

}

type Ready struct {
	Ready bool `json:"ready"`
}

func (h *handler) handleGetWarmer(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	isReader := new(Ready)

	db, err := mysql.Connect(appConfig.DBUsername, appConfig.DatabasePassword, appConfig.DBHost, appConfig.DBSchema)
	if err != nil {
		h.logger.WithError(err).Error("failed to create db connection")
		return apigw.RespondJSON(http.StatusInternalServerError, isReader, nil)
	}

	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		h.logger.WithError(err).Error("failed to ping db connection")
		return apigw.RespondJSON(http.StatusInternalServerError, isReader, nil)
	}

	isReader.Ready = true

	return apigw.RespondJSON(http.StatusOK, isReader, nil)

}
