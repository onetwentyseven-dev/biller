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
	logger    *logrus.Logger
	providers *mysql.ProvidersRepository
	bills     *mysql.BillsRepository
	gw        *apigw.Service
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

	providers := mysql.NewProviderRepository(db)
	bills := mysql.NewBillsRepository(db)

	api := apigw.New(logger)

	h := &handler{
		logger:    logger,
		providers: providers,
		bills:     bills,
		gw:        api,
	}

	api.AddHandler("GET /providers", h.handleGetProviders)
	api.AddHandler("POST /providers", h.handlePostProviders)
	api.AddHandler("GET /providers/{providerID}", h.handleGetProviderByID)
	api.AddHandler("PATCH /providers/{providerID}", h.handlePatchProviderByID)
	api.AddHandler("DELETE /providers/{providerID}", h.handleDeleteProviderByID)
	api.AddHandler("GET /providers/{providerID}/bills", h.handleGetBillsByProviderID)
	api.AddHandler("POST /providers/{providerID}/bills", h.handlePostBillsByProviderID)

	lambda.Start(api.HandleRoutes())

}

func (h *handler) handleGetProviders(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	providers, err := h.providers.Providers(ctx)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusBadRequest, "failed to query providers", nil, err)
	}

	return h.gw.RespondJSON(http.StatusOK, providers, nil)

}

func (h *handler) handleGetProviderByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	providerIDStr := event.PathParameters["providerID"]

	providerID, err := uuid.Parse(providerIDStr)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse provider id to valid uuid", nil, err)
	}

	provider, err := h.providers.Provider(ctx, providerID)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusBadRequest, "failed to query provider", nil, err)
	}

	return h.gw.RespondJSON(http.StatusOK, provider, nil)

}

func (h *handler) handlePostProviders(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	read := bytes.NewBufferString(event.Body)

	var provider = new(biller.Provider)
	err := json.NewDecoder(read).Decode(provider)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusBadRequest, "failed to decode request body", nil, err)
	}

	provider.ID = uuid.New()

	err = h.providers.CreateProvider(ctx, provider)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to create provider", nil, err)
	}

	return h.gw.RespondJSON(http.StatusOK, provider, nil)

}
func (h *handler) handlePatchProviderByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	providerIDStr := event.PathParameters["providerID"]

	providerID, err := uuid.Parse(providerIDStr)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse provider id to valid uuid", nil, err)
	}

	provider, err := h.providers.Provider(ctx, providerID)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch provider", nil, err)
	}

	read := bytes.NewBufferString(event.Body)

	err = json.NewDecoder(read).Decode(provider)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusBadRequest, "failed to decode request body", nil, err)
	}

	err = h.providers.UpdateProvider(ctx, providerID, provider)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to update provider", nil, err)
	}

	return h.gw.RespondJSON(http.StatusOK, provider, nil)

}

func (h *handler) handleDeleteProviderByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	providerIDStr := event.PathParameters["providerID"]

	providerID, err := uuid.Parse(providerIDStr)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse provider id to valid uuid", nil, err)
	}

	err = h.providers.DeleteProvider(ctx, providerID)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to delete provider", nil, err)
	}

	return h.gw.RespondJSON(http.StatusNoContent, nil, nil)

}

func (h *handler) handleGetBillsByProviderID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	providerIDStr := event.PathParameters["providerID"]

	providerID, err := uuid.Parse(providerIDStr)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse provider id to valid uuid", nil, err)
	}

	bills, err := h.bills.BillsByProvider(ctx, providerID)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusBadRequest, "failed to query bills", nil, err)
	}

	return h.gw.RespondJSON(http.StatusOK, bills, nil)

}

func (h *handler) handlePostBillsByProviderID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	providerIDStr := event.PathParameters["providerID"]

	providerID, err := uuid.Parse(providerIDStr)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse provider id to valid uuid", nil, err)
	}

	read := bytes.NewBufferString(event.Body)

	var bill = new(biller.Bill)
	err = json.NewDecoder(read).Decode(bill)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusBadRequest, "failed to decode request body", nil, err)
	}

	bill.ID = uuid.New()
	bill.ProviderID = providerID

	err = h.bills.CreateBill(ctx, bill)
	if err != nil {
		return h.gw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to create bill", nil, err)
	}

	return h.gw.RespondJSON(http.StatusOK, bill, nil)

}
