package voucher_handler

import (
	"customer-voucher-service/constants"
	"customer-voucher-service/constants/error_base"
	"customer-voucher-service/constants/message"
	pbVoucher "customer-voucher-service/protogen/voucher"
	"customer-voucher-service/services/voucher_service"
	"customer-voucher-service/utils/json_response"
	"fmt"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	voucherService voucher_service.IVoucherService
}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{voucherService: voucher_service.NewVoucherService()}
}

func VoucherRoutes(rg *gin.RouterGroup) {
	handler := NewHttpHandler()
	customer := rg.Group("/voucher")
	{
		customer.POST("/create", handler.CreateVoucher)
		customer.GET("/list", handler.ListVoucher)
		customer.GET("/detail", handler.DetailVoucher)
	}
}

func (h *HttpHandler) CreateVoucher(c *gin.Context) {
	payload := &pbVoucher.CreateVoucherReq{}
	if err := c.ShouldBindJSON(payload); err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrValidationFailed.HttpCode, error_base.ErrValidationFailed.Code, error_base.ErrValidationFailed.Message)
	}
	res, err := h.voucherService.CreateVoucher(c, payload)
	if err != nil && res != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrValidationFailed.HttpCode, error_base.ErrValidationFailed.Code, err.Error())
	}
	if err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrInternalServer.HttpCode, error_base.ErrInternalServer.Code, error_base.ErrInternalServer.Message)
	}
	json_response.Success(c, constants.CodeSystem, res)
}

func (h *HttpHandler) ListVoucher(c *gin.Context) {
	req := &pbVoucher.ListVoucherReq{}

	if brandIdStr := c.Query("brandId"); brandIdStr != "" {
		var brandId int32
		if _, err := fmt.Sscanf(brandIdStr, "%d", &brandId); err != nil {
			json_response.Error(c, constants.CodeSystem, error_base.ErrValidationFailed.HttpCode, error_base.ErrValidationFailed.Code, message.InvalidFormatMessage("brandId"))
		}
		req.BrandId = &brandId
	}

	res, err := h.voucherService.ListVoucher(c, req)
	if err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrInternalServer.HttpCode, error_base.ErrInternalServer.Code, error_base.ErrInternalServer.Message)
	}
	json_response.Success(c, constants.CodeSystem, res)
}

func (h *HttpHandler) DetailVoucher(c *gin.Context) {
	voucherIdStr := c.Query("voucherId")
	if voucherIdStr == "" {
		json_response.Error(c, constants.CodeSystem, error_base.ErrValidationFailed.HttpCode, error_base.ErrValidationFailed.Code, message.RequiredMessage("voucherId"))
	}
	var voucherId int32
	if _, err := fmt.Sscanf(voucherIdStr, "%d", &voucherId); err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrValidationFailed.HttpCode, error_base.ErrValidationFailed.Code, message.InvalidFormatMessage("voucherId"))
	}

	req := &pbVoucher.DetailVoucherReq{
		Id: voucherId,
	}

	res, err := h.voucherService.DetailVoucher(c, req)
	if err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrInternalServer.HttpCode, error_base.ErrInternalServer.Code, error_base.ErrInternalServer.Message)
	}
	json_response.Success(c, constants.CodeSystem, res)
}
