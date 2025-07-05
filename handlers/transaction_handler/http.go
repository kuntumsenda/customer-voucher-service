package transaction_handler

import (
	"customer-voucher-service/constants"
	"customer-voucher-service/constants/error_base"
	"customer-voucher-service/constants/message"
	pbTransaction "customer-voucher-service/protogen/transaction"
	"customer-voucher-service/services/transaction_service"
	"customer-voucher-service/utils/json_response"
	"fmt"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	transactionService transaction_service.ITransactionService
}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{transactionService: transaction_service.NewTransactionService()}
}

func TransactionRoutes(rg *gin.RouterGroup) {
	handler := NewHttpHandler()
	transaction := rg.Group("/transaction")
	{
		transaction.POST("/redemption", handler.TransactionRedeemPoint)
		transaction.GET("/list", handler.ListTransaction)
		transaction.GET("/detail", handler.DetailTransaction)
	}
}

func (h *HttpHandler) TransactionRedeemPoint(c *gin.Context) {
	payload := &pbTransaction.TransactionRedeemPointReq{}
	if err := c.ShouldBindJSON(payload); err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrValidationFailed.HttpCode, error_base.ErrValidationFailed.Code, error_base.ErrValidationFailed.Message)
	}
	res, err := h.transactionService.TransactionRedeemPoint(c, payload)
	if err != nil && res != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrValidationFailed.HttpCode, error_base.ErrValidationFailed.Code, err.Error())
	}
	if err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrInternalServer.HttpCode, error_base.ErrInternalServer.Code, error_base.ErrInternalServer.Message)
	}
	json_response.Success(c, constants.CodeSystem, res)
}

func (h *HttpHandler) ListTransaction(c *gin.Context) {
	req := &pbTransaction.ListTransactionReq{}

	if customerIdStr := c.Query("customerId"); customerIdStr != "" {
		var customerId int32
		if _, err := fmt.Sscanf(customerIdStr, "%d", &customerId); err != nil {
			json_response.Error(c, constants.CodeSystem, error_base.ErrValidationFailed.HttpCode, error_base.ErrValidationFailed.Code, "Invalid customerId format")
			return
		}
		req.CustomerId = &customerId
	}

	res, err := h.transactionService.ListTransaction(c, req)
	if err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrInternalServer.HttpCode, error_base.ErrInternalServer.Code, error_base.ErrInternalServer.Message)
	}
	json_response.Success(c, constants.CodeSystem, res)
}

func (h *HttpHandler) DetailTransaction(c *gin.Context) {
	transactionIdStr := c.Query("transactionId")
	if transactionIdStr == "" {
		json_response.Error(c, constants.CodeSystem, error_base.ErrValidationFailed.HttpCode, error_base.ErrValidationFailed.Code, message.RequiredMessage("transactionId"))
	}
	var transactionId int32
	if _, err := fmt.Sscanf(transactionIdStr, "%d", &transactionId); err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrValidationFailed.HttpCode, error_base.ErrValidationFailed.Code, "Invalid transactionId format")
		return
	}

	req := &pbTransaction.DetailTransactionReq{
		Id: transactionId,
	}

	res, err := h.transactionService.DetailTransaction(c, req)
	if err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrInternalServer.HttpCode, error_base.ErrInternalServer.Code, error_base.ErrInternalServer.Message)
		return
	}
	json_response.Success(c, constants.CodeSystem, res)
}
