package app

import (
	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
	"github.com/imraushankr/bervity/server/src/internal/handlers/v1"
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
	router := gin.Default()

	// Initialize core services
	authService := auth.NewAuth(&cfg.JWT)
	emailService := email.NewEmailService(&cfg.Email, log)
	storageProvider, err := storage.NewStorage(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB, log)
	authRepo := repository.NewAuthRepository(db.DB, log)

	// Initialize services
	authSvc := services.NewAuthService(
		authRepo,
		authService,
		emailService,
		cfg,
		log,
	)

	userSvc := services.NewUserService(
		userRepo,
		authService,
		emailService,
		cfg,
		storageProvider,
		log,
	)

	// Initialize handlers
	authHandler := v1.NewAuthHandler(authSvc, cfg, log)
	userHandler := v1.NewUserHandler(userSvc, log)
	healthHandler := v1.NewHealthHandler(cfg)

	// Setup routes with all required parameters
	routes.SetupRoutes(
		router, 
		userHandler, 
		healthHandler, 
		authHandler, 
		authService, 
		cfg,
		log, // Adding the missing logger parameter
	)

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	return router, nil
}