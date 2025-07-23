package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
	"github.com/imraushankr/bervity/server/src/internal/middleware"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
)

func RegisterAuthRoutes(r *gin.RouterGroup, h *v1.AuthHandler, authService *auth.Auth, cfg *configs.Config) {
	auth := r.Group("/auth")
	{
		// Public endpoints
		auth.POST("/signup", h.Register)
		auth.POST("/signin", h.Login)
		auth.POST("/signout", h.Logout)
		auth.GET("/verify-email", h.VerifyEmail)
		auth.POST("/forgot-password", h.InitiatePasswordReset)
		auth.POST("/password-reset", h.CompletePasswordReset)

		// Refresh token endpoint
		refreshGroup := auth.Group("/refresh")
		refreshGroup.Use(middleware.RefreshTokenAuth(authService, cfg))
		refreshGroup.POST("", h.RefreshToken)
	}
}