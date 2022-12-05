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
	logger    *logrus.Logger
	providers *mysql.ProvidersRepository
	bills     *mysql.BillsRepository
	gw        *apigw.Service
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

	providers := mysql.NewProviderRepository(db)
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
		logger:    logger,
		providers: providers,
		bills:     bills,
		gw:        api,
	}

	api.AddHandler(http.MethodGet, "/providers", h.handleGetProviders)
	api.AddHandler(http.MethodPost, "/providers", h.handlePostProviders)
	api.AddHandler(http.MethodGet, "/providers/{providerID}", h.handleGetProviderByID)
	api.AddHandler(http.MethodPatch, "/providers/{providerID}", h.handlePatchProviderByID)
	api.AddHandler(http.MethodDelete, "/providers/{providerID}", h.handleDeleteProviderByID)
	api.AddHandler(http.MethodGet, "/providers/{providerID}/bills", h.handleGetBillsByProviderID)
	api.AddHandler(http.MethodPost, "/providers/{providerID}/bills", h.handlePostBillsByProviderID)

	lambda.Start(apigw.UseMiddleware(api.HandleRoutes, authWare))

}

func (h *handler) handleGetProviders(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	providers, err := h.providers.Providers(ctx, userID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch providers", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, providers, nil)

}

func (h *handler) handleGetProviderByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	providerID, err := apigw.UUIDPathParameter("providerID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse provider id to valid uuid", nil, err)
	}

	provider, err := h.providers.Provider(ctx, userID, providerID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch provider", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, provider, nil)

}

func (h *handler) handlePostProviders(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	read := bytes.NewBufferString(event.Body)

	var provider = new(biller.Provider)
	err = json.NewDecoder(read).Decode(provider)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to decode request body", nil, err)
	}

	provider.ID = uuid.New()
	provider.UserID = userID

	err = h.providers.CreateProvider(ctx, provider)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to create provider", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, provider, nil)

}
func (h *handler) handlePatchProviderByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	providerID, err := apigw.UUIDPathParameter("providerID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse provider id to valid uuid", nil, err)
	}

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	provider, err := h.providers.Provider(ctx, userID, providerID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch provider", nil, err)
	}

	read := bytes.NewBufferString(event.Body)

	err = json.NewDecoder(read).Decode(provider)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to decode request body", nil, err)
	}

	err = h.providers.UpdateProvider(ctx, providerID, provider)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to update provider", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, provider, nil)

}

func (h *handler) handleDeleteProviderByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	providerID, err := apigw.UUIDPathParameter("providerID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse provider id to valid uuid", nil, err)
	}

	err = h.providers.DeleteProvider(ctx, providerID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to delete provider", nil, err)
	}

	return apigw.RespondJSON(http.StatusNoContent, nil, nil)

}

func (h *handler) handleGetBillsByProviderID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	providerID, err := apigw.UUIDPathParameter("providerID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse provider id to valid uuid", nil, err)
	}

	userID, err := apigw.UserIDFromContext(ctx)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to determine userID", nil, err)
	}

	bills, err := h.bills.BillsByProvider(ctx, userID, providerID)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch bills", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, bills, nil)

}

func (h *handler) handlePostBillsByProviderID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	providerID, err := apigw.UUIDPathParameter("providerID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse provider id to valid uuid", nil, err)
	}

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
	bill.ProviderID = providerID
	bill.UserID = userID

	err = h.bills.CreateBill(ctx, bill)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to create bill", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, bill, nil)

}
