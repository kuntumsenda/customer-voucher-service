package customer_handler

import (
	"customer-voucher-service/constants"
	"customer-voucher-service/constants/error_base"
	pbCustomer "customer-voucher-service/protogen/customer"
	"customer-voucher-service/services/customer_service"
	"customer-voucher-service/utils/json_response"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	customerService customer_service.ICustomerService
}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{customerService: customer_service.NewCustomerService()}
}

func CustomerRoutes(rg *gin.RouterGroup) {
	handler := NewHttpHandler()
	customer := rg.Group("/customer")
	{
		customer.POST("/create", handler.CreateCustomer)
		customer.GET("/list", handler.ListCustomer)
		customer.PUT("/update-points", handler.UpdateCustomerPoints)
	}
}

func (h *HttpHandler) CreateCustomer(c *gin.Context) {
	payload := &pbCustomer.CreateCustomerReq{}
	if err := c.ShouldBindJSON(payload); err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrValidationFailed.HttpCode, error_base.ErrValidationFailed.Code, error_base.ErrValidationFailed.Message)
	}
	res, err := h.customerService.CreateCustomer(c, payload)
	if err != nil && res != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrValidationFailed.HttpCode, error_base.ErrValidationFailed.Code, err.Error())
	}
	if err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrInternalServer.HttpCode, error_base.ErrInternalServer.Code, error_base.ErrInternalServer.Message)
	}
	json_response.Success(c, constants.CodeSystem, res)
}

func (h *HttpHandler) ListCustomer(c *gin.Context) {
	req := &pbCustomer.ListCustomerReq{}
	res, err := h.customerService.ListCustomer(c, req)
	if err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrInternalServer.HttpCode, error_base.ErrInternalServer.Code, error_base.ErrInternalServer.Message)
	}
	json_response.Success(c, constants.CodeSystem, res)
}

func (h *HttpHandler) UpdateCustomerPoints(c *gin.Context) {
	payload := &pbCustomer.UpdateCustomerPointsReq{}
	if err := c.ShouldBindJSON(payload); err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrValidationFailed.HttpCode, error_base.ErrValidationFailed.Code, error_base.ErrValidationFailed.Message)
	}
	res, err := h.customerService.UpdateCustomerPoints(c, payload)
	if err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrInternalServer.HttpCode, error_base.ErrInternalServer.Code, error_base.ErrInternalServer.Message)
	}
	json_response.Success(c, constants.CodeSystem, res)
}
