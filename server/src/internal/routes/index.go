package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	routerv1 "github.com/imraushankr/bervity/server/src/internal/routes/v1"
)

func SetupRoutes(
	router *gin.Engine, 
	userHandler *v1.UserHandler,
	healthHandler *v1.HealthHandler, 
	authHandler *v1.AuthHandler,
	urlHandler *v1.URLHandler,
	creditHandler *v1.CreditHandler,
	subHandler *v1.SubscriptionHandler,
	authService *auth.Auth, 
	urlRepo interfaces.URLRepository,
	cfg *configs.Config, 
	log logger.Logger,
) {
	api := router.Group("/api")
	{
		v1Group := api.Group("/v1")
		routerv1.RegisterAuthRoutes(v1Group, authHandler, authService, cfg, log)
		routerv1.RegisterUserRoutes(v1Group, userHandler, authService, cfg, log)
		routerv1.RegisterURLRoutes(v1Group, urlHandler, authService, urlRepo, cfg, log)
		routerv1.RegisterCreditRoutes(v1Group, creditHandler, authService, cfg, log)
		routerv1.RegisterSubscriptionRoutes(v1Group, subHandler, authService, cfg, log)
		routerv1.RegisterSystemRoutes(v1Group, healthHandler)
	}
}