// package v1

// import (
// 	"github.com/gin-gonic/gin"
// 	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
// 	"github.com/imraushankr/bervity/server/src/internal/middleware"
// )

// func RegisterAuthRoutes(r *gin.RouterGroup, h *v1.UserHandler) {
// 	auth := r.Group("/auth")
// 	{
// 		// Public endpoints
// 		auth.POST("/signup", h.Register)
// 		auth.POST("/signin", h.Login)
// 		auth.GET("/verify-email", h.VerifyEmail)
// 		auth.POST("/password-reset", h.InitiatePasswordReset)
// 		auth.POST("/password-reset/confirm", h.CompletePasswordReset)

// 		// Refresh token endpoints
// 		refreshGroup := auth.Group("", middleware.RefreshTokenAuth(auth))
// 	}
// }

package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
	"github.com/imraushankr/bervity/server/src/internal/middleware"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
)

func RegisterAuthRoutes(r *gin.RouterGroup, h *v1.UserHandler, authService *auth.Auth, cfg *configs.Config) {
	auth := r.Group("/auth")
	{
		// Public endpoints
		auth.POST("/signup", h.Register)
		auth.POST("/signin", h.Login)
		auth.GET("/verify-email", h.VerifyEmail)
		auth.POST("/password-reset", h.InitiatePasswordReset)
		auth.POST("/password-reset/confirm", h.CompletePasswordReset)

		// Refresh token endpoint
		refreshGroup := auth.Group("/refresh")
		refreshGroup.Use(middleware.RefreshTokenAuth(authService, cfg))
		refreshGroup.POST("", h.RefreshToken)
	}
}