package app

import (
	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
	"github.com/imraushankr/bervity/server/src/internal/middleware"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	"github.com/imraushankr/bervity/server/src/internal/pkg/database"
	"github.com/imraushankr/bervity/server/src/internal/pkg/email"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"github.com/imraushankr/bervity/server/src/internal/pkg/storage"
	"github.com/imraushankr/bervity/server/src/internal/repository"
	"github.com/imraushankr/bervity/server/src/internal/routes"
	"github.com/imraushankr/bervity/server/src/internal/services"
)

func SetupRouter(cfg *configs.Config, db *database.DB, log logger.Logger) (*gin.Engine, error) {
	// Set Gin mode
	if cfg.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logs(log))
	router.Use(middleware.CORSMiddleware(cfg))

	// Initialize dependencies
	authService := auth.NewAuth(&cfg.JWT)
	emailService := email.NewEmailService(&cfg.Email, log)
	storageService, err := storage.NewStorage(cfg)
	if err != nil {
		return nil, err
	}

	// Since DB embeds *gorm.DB, we can use it directly
	userRepo := repository.NewUserRepository(db.DB, log)

	// Initialize services
	userService := services.NewUserService(
		userRepo,
		authService,
		emailService,
		cfg,
		storageService,
		log,
	)

	// Initialize handlers
	userHandler := v1.NewUserHandler(userService, log)

	// Setup routes
	routes.SetupRoutes(router, userHandler)
	return router, nil
}
