// package routes

// import (
// 	"github.com/gin-gonic/gin"
// 	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
// 	routerv1 "github.com/imraushankr/bervity/server/src/internal/routes/v1"
// )

// func SetupRoutes(router *gin.Engine, userHandler *v1.UserHandler) {
// 	api := router.Group("/api")
// 	{
// 		routerv1.AuthRoutes(api.Group("/v1"), userHandler)
// 		routerv1.UserRoutes(api.Group("/v1"), userHandler)
// 		routerv1.RegisterSystemRoutes(api.Group("/v1"), v1.NewHealthHandler())
// 	}
// }


// package routes

// import (
// 	"github.com/gin-gonic/gin"
// 	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
// 	routerv1 "github.com/imraushankr/bervity/server/src/internal/routes/v1"
// )

// func SetupRoutes(router *gin.Engine, userHandler *v1.UserHandler, healthHandler *v1.HealthHandler) {
// 	api := router.Group("/api")
// 	{
// 		routerv1.RegisterAuthRoutes(api.Group("/v1"), userHandler)
// 		routerv1.RegisterUserRoutes(api.Group("/v1"), userHandler)
// 		routerv1.RegisterSystemRoutes(api.Group("/v1"), healthHandler)
// 	}
// }


package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	routerv1 "github.com/imraushankr/bervity/server/src/internal/routes/v1"
)

func SetupRoutes(router *gin.Engine, userHandler *v1.UserHandler, healthHandler *v1.HealthHandler, authService *auth.Auth, cfg *configs.Config) {
	api := router.Group("/api")
	{
		routerv1.RegisterAuthRoutes(api.Group("/v1"), userHandler, authService, cfg)
		routerv1.RegisterUserRoutes(api.Group("/v1"), userHandler)
		routerv1.RegisterSystemRoutes(api.Group("/v1"), healthHandler)
	}
}