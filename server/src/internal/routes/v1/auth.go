package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
	"github.com/imraushankr/bervity/server/src/internal/middleware"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

func RegisterAuthRoutes(r *gin.RouterGroup, h *v1.AuthHandler, auth *auth.Auth, cfg *configs.Config, log logger.Logger) {
	authGroup := r.Group("/auth")
	{
		// Public endpoints
		authGroup.POST("/signup", h.Register)
		authGroup.POST("/signin", h.Login)
		authGroup.POST("/signout", h.Logout)
		authGroup.GET("/verify-email", h.VerifyEmail)
		authGroup.POST("/forgot-password", h.InitiatePasswordReset)
		authGroup.PATCH("/reset-password/:token", h.CompletePasswordReset)

		// Protected endpoints
		protected := authGroup.Group("")
		protected.Use(middleware.JWTAuth(auth, cfg, log))
		{
			protected.PATCH("/change-password", h.ChangePassword)
		}

		// Refresh token endpoint
		refreshGroup := authGroup.Group("/refresh")
		refreshGroup.Use(middleware.RefreshTokenAuth(auth, cfg, log))
		refreshGroup.POST("", h.RefreshToken)
	}
}