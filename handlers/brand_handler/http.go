package brand_handler

import (
	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	authService auth_service.IAuthService
}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{authService: brand_service.NewAuthService()}
}

func RegisterAuthRoutes(rgroup *gin.RouterGroup) {
	auth := rgroup.Group("/auth")
	{
		auth.POST("/login", handler.Login)
	}
}
