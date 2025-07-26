package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
	"github.com/imraushankr/bervity/server/src/internal/middleware"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

func RegisterUserRoutes(r *gin.RouterGroup, h *v1.UserHandler, auth *auth.Auth, cfg *configs.Config, log logger.Logger) {
	users := r.Group("/users")
	users.Use(middleware.JWTAuth(auth, cfg, log))
	{
		// Profile management
		users.GET("/me", h.GetProfile)
		users.PUT("/me", h.UpdateProfile)

		// Avatar
		users.POST("/avatar", h.UploadAvatar)

		// Account management
		users.DELETE("/me", h.DeleteAccount)
	}
}
