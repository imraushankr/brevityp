package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
	"github.com/imraushankr/bervity/server/src/internal/middleware"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

func RegisterCreditRoutes(router *gin.RouterGroup, creditHandler *v1.CreditHandler, authService *auth.Auth, cfg *configs.Config, log logger.Logger) {
	creditRoutes := router.Group("/credits")
	{
		creditRoutes.Use(middleware.JWTAuth(authService, cfg, log))

		creditRoutes.GET("/balance", creditHandler.GetBalance)
		creditRoutes.POST("/apply-promo", creditHandler.ApplyPromoCode)
		creditRoutes.GET("/usage", creditHandler.GetUsage)
	}
}