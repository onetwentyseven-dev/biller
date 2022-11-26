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
	bills  *mysql.BillsRepository
	gw     *apigw.Service
}

func main() {

	awsCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("failed to initiali aws config")
	}

	loadConfig(awsCfg)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	db, err := mysql.Connect(appConfig.DBUsername, appConfig.DatabasePassword, appConfig.DBHost, appConfig.DBSchema)
	if err != nil {
		logger.WithError(err).Fatal("failed to connect to db")
	}

	bills := mysql.NewBillsRepository(db)

	api := apigw.New(logger)

	h := &handler{
		logger: logger,
		bills:  bills,
		gw:     api,
	}

	api.AddHandler("GET /bills", h.handleGetBills)

	lambda.Start(api.HandleRoutes())

}

func (h *handler) handleGetBills(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	bills, err := h.bills.Bills(ctx)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusBadRequest, "failed to query bills", nil, err)
	}

	return h.gw.RespondJSON(http.StatusOK, bills, nil)

}
