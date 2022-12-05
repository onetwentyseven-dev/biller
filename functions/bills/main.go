package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/google/uuid"
	"github.com/onetwentyseven-dev/apigw"
	"github.com/onetwentyseven-dev/biller"
	"github.com/onetwentyseven-dev/biller/internal/mysql"
	"github.com/sirupsen/logrus"
)

type handler struct {
	logger *logrus.Logger
	bills  *mysql.BillsRepository
}

func main() {

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	awsCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize aws config")
	}

	loadConfig(awsCfg)

	db, err := mysql.Connect(appConfig.DBUsername, appConfig.DatabasePassword, appConfig.DBHost, appConfig.DBSchema)
	if err != nil {
		logger.WithError(err).Fatal("failed to connect to db")
	}

	bills := mysql.NewBillsRepository(db)

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	authWare, err := apigw.Auth(client, appConfig.AuthTenant, appConfig.AuthClientID, appConfig.AuthAudience)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize auth middleware")
	}

	api := apigw.New(logger)

	h := &handler{
		logger: logger,
		bills:  bills,
	}

	api.AddHandler(http.MethodGet, "/bills", h.handleGetBills)
	api.AddHandler(http.MethodPost, "/bills", h.handlePostBills)
	api.AddHandler(http.MethodGet, "/bills/{billID}", h.handleGetBillByID)
	api.AddHandler(http.MethodPatch, "/bills/{billID}", h.handlePatchBillByID)
	api.AddHandler(http.MethodDelete, "/bills/{billID}", h.handleDeleteBillByID)

	lambda.Start(apigw.UseMiddleware(api.HandleRoutes, authWare))

}

func (h *handler) handleGetBills(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	bills, err := h.bills.Bills(ctx, userID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch bills", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, bills, nil)

}

func (h *handler) handleGetBillByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	billID, err := apigw.UUIDPathParameter("billID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse bill id to valid uuid", nil, err)
	}

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	bill, err := h.bills.Bill(ctx, userID, billID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch bill", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, bill, nil)

}

func (h *handler) handlePostBills(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	read := bytes.NewBufferString(event.Body)

	var bill = new(biller.Bill)
	err = json.NewDecoder(read).Decode(bill)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to decode request body", nil, err)
	}

	bill.ID = uuid.New()
	bill.UserID = userID

	err = h.bills.CreateBill(ctx, bill)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to create bill", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, bill, nil)

}

func (h *handler) handlePatchBillByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	billID, err := apigw.UUIDPathParameter("billID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse bill id to valid uuid", nil, err)
	}

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	bill, err := h.bills.Bill(ctx, userID, billID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch bill", nil, err)
	}

	read := bytes.NewBufferString(event.Body)

	err = json.NewDecoder(read).Decode(bill)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to decode request body", nil, err)
	}

	err = h.bills.UpdateBill(ctx, billID, bill)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to update bill", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, bill, nil)

}

func (h *handler) handleDeleteBillByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	billID, err := apigw.UUIDPathParameter("billID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse bill id to valid uuid", nil, err)
	}

	err = h.bills.DeleteBill(ctx, billID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to delete bill", nil, err)
	}

	return apigw.RespondJSON(http.StatusNoContent, nil, nil)

}
