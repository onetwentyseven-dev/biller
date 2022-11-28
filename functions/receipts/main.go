package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/onetwentyseven-dev/apigw"
	"github.com/onetwentyseven-dev/biller"
	"github.com/onetwentyseven-dev/biller/internal/mysql"
	"github.com/sirupsen/logrus"
)

type handler struct {
	logger   *logrus.Logger
	receipts *mysql.ReceiptRepository
	s3       *s3.Client
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

	receipts := mysql.NewReceiptRepository(db)

	api := apigw.New(logger)

	h := &handler{
		logger:   logger,
		receipts: receipts,
		s3:       s3.NewFromConfig(awsCfg),
	}

	api.AddHandler(http.MethodGet, "/receipts", h.handleGetReceipts)
	api.AddHandler(http.MethodPost, "/receipts", h.handlePostReceipt)
	api.AddHandler(http.MethodGet, "/receipts/{receiptID}", h.handleGetReceipt)
	api.AddHandler(http.MethodGet, "/receipts/{receiptID}/file", h.handleGetReceiptFile)
	api.AddHandler(http.MethodPost, "/receipts/{receiptID}/file", h.handlePostReceiptFile)

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

	read := bytes.NewBufferString(event.Body)

	var receipt = new(biller.Receipt)
	err := json.NewDecoder(read).Decode(receipt)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to decode request body", nil, err)
	}

	receipt.ID = uuid.New()

	err = h.receipts.CreateReceipt(ctx, receipt)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to create receipt", nil, err)
	}

	return apigw.RespondJSON(http.StatusOK, receipt, nil)
}

func (h *handler) handleGetReceiptFile(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	receiptID, err := apigw.UUIDPathParameter("receiptID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse receipt id to valid uuid", nil, err)
	}

	getOutput, err := h.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(appConfig.ReceiptBucket),
		Key:    aws.String(receiptID.String()),
	})
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to fetch object from bucket", nil, err)
	}

	buf := new(bytes.Buffer)

	encoder := base64.NewEncoder(base64.StdEncoding, buf)

	n, err := io.Copy(encoder, getOutput.Body)
	encoder.Close()
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "fatal error read s3 object", nil, err)
	}

	return apigw.Respond(http.StatusOK, buf.String(), map[string]string{
		"content-length": strconv.FormatInt(n, 10),
		"content-type":   aws.ToString(getOutput.ContentType),
	}, true)

}

func (h *handler) handlePostReceiptFile(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	receiptID, err := apigw.UUIDPathParameter("receiptID", &event)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to parse receipt id to valid uuid", nil, err)
	}

	contentType := event.Headers["content-type"]
	if contentType == "" {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "required header content-type missing from request", nil, err)
	}

	bodyBytes, err := base64.StdEncoding.DecodeString(event.Body)
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusBadRequest, "failed to base64 decode file", nil, err)
	}

	read := bytes.NewBuffer(bodyBytes)

	_, err = h.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(appConfig.ReceiptBucket),
		Key:         aws.String(receiptID.String()),
		ContentType: aws.String(contentType),
		Body:        read,
	})
	if err != nil {
		return apigw.RespondJSONError(ctx, http.StatusInternalServerError, "failed to put file as object into S3", nil, err)
	}

	return apigw.RespondJSON(http.StatusNoContent, nil, nil)

}
