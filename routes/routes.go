package routes

import (
	"customer-voucher-service/handlers/brand_handler"
	"github.com/gin-gonic/gin"
)

func ApiRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		brand_handler.BrandRoutes(api)
	}
}
