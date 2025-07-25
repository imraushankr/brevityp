// package v1

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/imraushankr/bervity/server/src/configs"
// 	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
// 	"github.com/imraushankr/bervity/server/src/internal/middleware"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
// )

// func RegisterURLRoutes(router *gin.RouterGroup, urlHandler *v1.URLHandler, authService *auth.Auth, cfg *configs.Config, log logger.Logger) {
// 	urlRoutes := router.Group("/urls")
// 	{
// 		urlRoutes.Use(middleware.JWTAuth(authService, cfg, log))

// 		urlRoutes.POST("", urlHandler.CreateURL)
// 		urlRoutes.GET("", urlHandler.GetUserURLs)
// 		urlRoutes.GET("/:code", urlHandler.GetURL)
// 		urlRoutes.PUT("/:id", urlHandler.UpdateURL)
// 		urlRoutes.DELETE("/:id", urlHandler.DeleteURL)
// 		urlRoutes.GET("/:id/analytics", urlHandler.GetAnalytics)
// 	}

// 	router.GET("/:code", urlHandler.Redirect)
// }


package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
	"github.com/imraushankr/bervity/server/src/internal/middleware"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

func RegisterURLRoutes(router *gin.RouterGroup, urlHandler *v1.URLHandler, authService *auth.Auth, cfg *configs.Config, log logger.Logger) {
	urlRoutes := router.Group("/urls")
	{
		urlRoutes.Use(middleware.JWTAuth(authService, cfg, log))

		urlRoutes.POST("", urlHandler.CreateURL)
		urlRoutes.GET("", urlHandler.GetUserURLs)
		urlRoutes.GET("/:id", urlHandler.GetURL) 
		urlRoutes.PUT("/:id", urlHandler.UpdateURL)
		urlRoutes.DELETE("/:id", urlHandler.DeleteURL)
		urlRoutes.GET("/:id/analytics", urlHandler.GetAnalytics)
	}

	// Keep the redirect route with :code since it's public and needs the short code
	router.GET("/:code", urlHandler.Redirect)
}