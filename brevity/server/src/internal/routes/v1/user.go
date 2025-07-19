package v1

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
)

func UserRoutes(router *gin.RouterGroup, userHandler *v1.UserHandler) {
	users := router.Group("/users")
	{
		users.GET("/u", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "User endpoint"})
		})
	}
	// users.Use()
}
