package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

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
	sheets *mysql.BillSheetRepository
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

	sheets := mysql.NewBillSheetRepository(db)

	api := apigw.New(logger)

	h := &handler{
		logger: logger,
		sheets: sheets,
	}

	api.AddHandlerMethod(http.MethodGet, "/sheets", h.handleGetSheets)
	api.AddHandlerMethod(http.MethodPost, "/sheets", h.handlePostSheets)
	api.AddHandlerMethod(http.MethodGet, "/sheets/{sheetID}", h.handleGetSheet)
	api.AddHandlerMethod(http.MethodPatch, "/sheets/{sheetID}", h.handlePatchSheetByID)
	api.AddHandlerMethod(http.MethodDelete, "/sheets/{sheetID}", h.handleDeleteSheetByID)

	lambda.Start(api.HandleRoutes)

}

func (h *handler) handleGetSheets(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	sheets, err := h.sheets.Sheets(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch sheets", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, sheets, nil)

}

func (h *handler) handlePostSheets(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	read := bytes.NewBufferString(event.Body)

	var sheet = new(biller.BillSheet)
	err := json.NewDecoder(read).Decode(sheet)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to decode request body", nil, err)
	}

	sheet.ID = uuid.New()

	err = h.sheets.CreateSheet(ctx, sheet)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to create bill", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, sheet, nil)

}

func (h *handler) handleGetSheet(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	sheetIDStr := event.PathParameters["sheetID"]

	sheetID, err := uuid.Parse(sheetIDStr)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse bill id to valid uuid", nil, err)
	}

	sheet, err := h.sheets.Sheet(ctx, sheetID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch sheets", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, sheet, nil)

}

func (h *handler) handlePatchSheetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	sheetIDStr := event.PathParameters["sheetID"]

	sheetID, err := uuid.Parse(sheetIDStr)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse bill id to valid uuid", nil, err)
	}

	sheet, err := h.sheets.Sheet(ctx, sheetID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch sheets", nil, err)
	}

	read := bytes.NewBufferString(event.Body)

	err = json.NewDecoder(read).Decode(sheet)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to decode request body", nil, err)
	}

	err = h.sheets.UpdateSheet(ctx, sheetID, sheet)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to update sheet", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, sheet, nil)

}

func (h *handler) handleDeleteSheetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	sheetIDStr := event.PathParameters["sheetID"]

	sheetID, err := uuid.Parse(sheetIDStr)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse bill id to valid uuid", nil, err)
	}

	err = h.sheets.DeleteSheet(ctx, sheetID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to delete sheet", nil, err)
	}

	return apigw.RespondJSON(http.StatusNoContent, nil, nil)

}
