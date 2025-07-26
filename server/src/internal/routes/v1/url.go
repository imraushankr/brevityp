package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
	"github.com/imraushankr/bervity/server/src/internal/middleware"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

func RegisterURLRoutes(
	router *gin.RouterGroup,
	urlHandler *v1.URLHandler,
	authService *auth.Auth,
	urlRepo interfaces.URLRepository,
	cfg *configs.Config,
	log logger.Logger,
) {
	// Public routes (no auth required)
	router.POST("/urls",
		middleware.AnonymousURLLimit(urlRepo, log, cfg.App.AnonURLLimit),
		urlHandler.CreateURL,
	)
	// Add the API prefix to the redirect route
	router.GET("/r/:code", urlHandler.Redirect)

	// Authenticated routes
	authRoutes := router.Group("/urls")
	authRoutes.Use(middleware.JWTAuth(authService, cfg, log))
	{
		authRoutes.GET("", urlHandler.GetUserURLs)
		authRoutes.GET("/:id", urlHandler.GetURL)
		authRoutes.PUT("/:id", urlHandler.UpdateURL)
		authRoutes.DELETE("/:id", urlHandler.DeleteURL)
		authRoutes.GET("/:id/analytics", urlHandler.GetAnalytics)
	}
}