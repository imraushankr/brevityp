package routes

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
	routerv1 "github.com/imraushankr/bervity/server/src/internal/routes/v1"
)

func SetupRoutes(router *gin.Engine, userHandler *v1.UserHandler) {
	api := router.Group("/api")
	{
		routerv1.AuthRoutes(api.Group("/v1"), userHandler)
		routerv1.UserRoutes(api.Group("/v1"), userHandler)
	}
}
