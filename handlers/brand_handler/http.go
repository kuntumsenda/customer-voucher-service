package brand_handler

import (
	"customer-voucher-service/constants"
	"customer-voucher-service/constants/error_base"
	pbBrand "customer-voucher-service/protogen/brand"
	"customer-voucher-service/services/brand_service"
	"customer-voucher-service/utils/json_response"
	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	brandService brand_service.IBrandService
}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{brandService: brand_service.NewBrandService()}
}

func BrandRoutes(rg *gin.RouterGroup) {
	handler := NewHttpHandler()
	brand := rg.Group("/brand")
	{
		brand.POST("/create", handler.CreateBrand)
	}
}

func (h *HttpHandler) CreateBrand(c *gin.Context) {
	payload := &pbBrand.CreateBrandReq{}
	res, err := h.brandService.CreateBrand(c, payload)
	if err != nil {
		json_response.Error(c, constants.CodeSystem, error_base.ErrInternalServer.HttpCode, error_base.ErrInternalServer.Code, error_base.ErrInternalServer.Message)
	}
	json_response.Success(c, constants.CodeSystem, res)
}
