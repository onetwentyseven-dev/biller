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
	logger   *logrus.Logger
	sheets   *mysql.BillSheetRepository
	bills    *mysql.BillsRepository
	receipts *mysql.ReceiptRepository
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

	sheets := mysql.NewBillSheetRepository(db)
	bills := mysql.NewBillsRepository(db)
	receipts := mysql.NewReceiptRepository(db)

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	authWare, err := apigw.Auth(client, appConfig.AuthTenant, appConfig.AuthClientID, appConfig.AuthAudience)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize auth middleware")
	}

	api := apigw.New(logger)

	h := &handler{
		logger:   logger,
		sheets:   sheets,
		bills:    bills,
		receipts: receipts,
	}

	api.AddHandler(http.MethodGet, "/sheets", h.handleGetSheets)
	api.AddHandler(http.MethodPost, "/sheets", h.handlePostSheets)
	api.AddHandler(http.MethodGet, "/sheets/{sheetID}", h.handleGetSheet)
	api.AddHandler(http.MethodPatch, "/sheets/{sheetID}", h.handlePatchSheetByID)
	api.AddHandler(http.MethodDelete, "/sheets/{sheetID}", h.handleDeleteSheetByID)
	api.AddHandler(http.MethodGet, "/sheets/{sheetID}/entries", h.handleGetSheetEntries)
	api.AddHandler(http.MethodPost, "/sheets/{sheetID}/entries", h.handlePostSheetEntries)
	api.AddHandler(http.MethodPatch, "/sheets/{sheetID}/entries/{entryID}", h.handlePatchSheetEntry)
	api.AddHandler(http.MethodDelete, "/sheets/{sheetID}/entries/{entryID}", h.handleDeleteSheetEntry)

	lambda.Start(apigw.UseMiddleware(api.HandleRoutes, authWare))

}

func (h *handler) handleGetSheets(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	sheets, err := h.sheets.Sheets(ctx, userID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch sheets", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, sheets, nil)

}

func (h *handler) handlePostSheets(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	read := bytes.NewBufferString(event.Body)

	var sheet = new(biller.BillSheet)
	err = json.NewDecoder(read).Decode(sheet)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to decode request body", nil, err)
	}

	sheet.ID = uuid.New()
	sheet.UserID = userID

	err = h.sheets.CreateSheet(ctx, sheet)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to create bill", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, sheet, nil)

}

func (h *handler) handleGetSheet(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	sheetID, err := apigw.UUIDPathParameter("sheetID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse sheet id to valid uuid", nil, err)
	}

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	sheet, err := h.sheets.Sheet(ctx, userID, sheetID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch sheets", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, sheet, nil)

}

func (h *handler) handlePatchSheetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	sheetID, err := apigw.UUIDPathParameter("sheetID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse sheet id to valid uuid", nil, err)
	}

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	sheet, err := h.sheets.Sheet(ctx, userID, sheetID)
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

	sheetID, err := apigw.UUIDPathParameter("sheetID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse sheet id to valid uuid", nil, err)
	}

	err = h.sheets.DeleteSheet(ctx, sheetID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to delete sheet", nil, err)
	}

	return apigw.RespondJSON(http.StatusNoContent, nil, nil)

}

func (h *handler) handleGetSheetEntries(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	sheetID, err := apigw.UUIDPathParameter("sheetID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse sheet id to valid uuid", nil, err)
	}

	entries, err := h.sheets.SheetEntries(ctx, sheetID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch sheet entries", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, entries, nil)
}

func (h *handler) handlePostSheetEntries(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	sheetID, err := apigw.UUIDPathParameter("sheetID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse sheet id to valid uuid", nil, err)
	}

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	_, err = h.sheets.Sheet(ctx, userID, sheetID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch sheet", nil, err)
	}

	read := bytes.NewBufferString(event.Body)

	var entry = new(biller.BillSheetEntry)
	err = json.NewDecoder(read).Decode(entry)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to decode request body", nil, err)
	}

	_, err = h.bills.Bill(ctx, userID, entry.BillID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch bill", nil, err)
	}

	if entry.ReceiptID != nil {
		_, err = h.receipts.Receipt(ctx, userID, *entry.ReceiptID)
		if err != nil {
			return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch receipt", nil, err)
		}
	}

	entry.EntryID = uuid.New()
	entry.SheetID = sheetID

	err = h.sheets.CreateSheetEntry(ctx, entry)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to create entry", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, entry, nil)

}

func (h *handler) handlePatchSheetEntry(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	sheetID, err := apigw.UUIDPathParameter("sheetID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse sheet id to valid uuid", nil, err)
	}

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	_, err = h.sheets.Sheet(ctx, userID, sheetID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch sheet", nil, err)
	}

	entryID, err := apigw.UUIDPathParameter("entryID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse entry id to valid uuid", nil, err)
	}

	entry, err := h.sheets.SheetEntry(ctx, sheetID, entryID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch sheet for provider sheetID and entryID", nil, err)
	}

	read := bytes.NewBufferString(event.Body)

	err = json.NewDecoder(read).Decode(entry)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to decode request body", nil, err)
	}

	_, err = h.bills.Bill(ctx, userID, entry.BillID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch bill", nil, err)
	}

	if entry.ReceiptID != nil {
		_, err = h.receipts.Receipt(ctx, userID, *entry.ReceiptID)
		if err != nil {
			return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch receipt", nil, err)
		}
	}

	entry.EntryID = entryID
	entry.SheetID = sheetID

	err = h.sheets.UpdateSheetEntry(ctx, entryID, entry)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to update entry", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, entry, nil)

}

func (h *handler) handleDeleteSheetEntry(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	sheetID, err := apigw.UUIDPathParameter("sheetID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse sheet id to valid uuid", nil, err)
	}

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	_, err = h.sheets.Sheet(ctx, userID, sheetID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch sheet", nil, err)
	}

	entryID, err := apigw.UUIDPathParameter("entryID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse entry id to valid uuid", nil, err)
	}

	err = h.sheets.DeleteSheetEntry(ctx, sheetID, entryID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to delete entry", nil, err)
	}

	return apigw.RespondJSON(http.StatusNoContent, nil, nil)

}
