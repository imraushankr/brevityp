// routers/v1/user.go
package v1

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
)

func RegisterUserRoutes(r *gin.RouterGroup, h *v1.UserHandler) {
	users := r.Group("/users")
	// users.Use(middleware.JWTAuth()) // Require authentication

	{
		// // Profile management
		// users.GET("/me", h.GetProfile)
		// users.PUT("/me", h.UpdateProfile)
		// users.PATCH("/me", h.PartialUpdateProfile)

		// // Security
		// users.PUT("/password", h.ChangePassword)
		// users.POST("/logout", h.LogoutAllSessions)

		// // Avatar
		// users.POST("/avatar", h.UploadAvatar)
		// users.DELETE("/avatar", h.DeleteAvatar)

		// // Account management
		// users.DELETE("/account", h.DeleteAccount)

		users.GET("/me", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "This is a placeholder for user profile retrieval",
			})
		})
	}
}