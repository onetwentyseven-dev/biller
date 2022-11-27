package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/davecgh/go-spew/spew"
	"github.com/onetwentyseven-dev/apigw"
	"github.com/onetwentyseven-dev/biller/internal/mysql"
	"github.com/sirupsen/logrus"
)

type handler struct {
	logger   *logrus.Logger
	receipts *mysql.ReceiptRepository
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

	receipts := mysql.NewReceiptRepository(db)

	api := apigw.New(logger)

	h := &handler{
		logger:   logger,
		receipts: receipts,
	}

	api.AddHandler(http.MethodGet, "/receipts", h.handleGetReceipts)
	api.AddHandler(http.MethodPost, "/receipts", h.handlePostReceipt)
	api.AddHandler(http.MethodGet, "/receipts/{receiptID}", h.handleGetReceipt)

	lambda.Start(api.HandleRoutes)

}

func (h *handler) handleGetReceipts(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	receipts, err := h.receipts.Receipts(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch receipts", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, receipts, nil)

}

func (h *handler) handleGetReceipt(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	receiptID, err := apigw.UUIDPathParameter("receiptID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse receipt id to valid uuid", nil, err)
	}

	receipt, err := h.receipts.Receipt(ctx, receiptID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch receipt", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, receipt, nil)

}

func (h *handler) handlePostReceipt(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	spew.Dump(event)
	return apigw.RespondJSON(http.StatusNoContent, nil, nil)
}
