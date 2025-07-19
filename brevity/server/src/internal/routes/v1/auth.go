package v1

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
)

func AuthRoutes(router *gin.RouterGroup, userHandler *v1.UserHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", userHandler.Register)
	}
}
