package routes

import (
	"customer-voucher-service/handlers/brand_handler"
	"customer-voucher-service/handlers/customer_handler"
	"customer-voucher-service/handlers/transaction_handler"
	"customer-voucher-service/handlers/voucher_handler"

	"github.com/gin-gonic/gin"
)

func ApiRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		brand_handler.BrandRoutes(api)
		customer_handler.CustomerRoutes(api)
		voucher_handler.VoucherRoutes(api)
		transaction_handler.TransactionRoutes(api)
	}
}
