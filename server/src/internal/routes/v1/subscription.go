package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
	"github.com/imraushankr/bervity/server/src/internal/middleware"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

func RegisterSubscriptionRoutes(router *gin.RouterGroup, subHandler *v1.SubscriptionHandler, authService *auth.Auth, cfg *configs.Config, log logger.Logger) {
	subRoutes := router.Group("/subscriptions")
	{
		subRoutes.Use(middleware.JWTAuth(authService, cfg, log))

		subRoutes.POST("", subHandler.CreateSubscription)
		subRoutes.GET("", subHandler.GetSubscription)
		subRoutes.PUT("", subHandler.UpdateSubscription)
		subRoutes.DELETE("", subHandler.CancelSubscription)
		subRoutes.GET("/plans", subHandler.GetPlans)
		subRoutes.GET("/payments", subHandler.GetPaymentHistory)
	}
}