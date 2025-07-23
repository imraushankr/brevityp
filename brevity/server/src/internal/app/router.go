// // package app

// // import (
// // 	"github.com/gin-gonic/gin"
// // 	"github.com/imraushankr/bervity/server/src/configs"
// // 	"github.com/imraushankr/bervity/server/src/internal/pkg/database"
// // 	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
// // )

// // type RouterDependencies struct {
// // }

// // func SetupRouter(cfg *configs.Config, db *database.DB, log logger.Logger) (*gin.Engine, error) {
// // 	// Initialize auth service
// // 	// Initialize email service
// // 	// Initialize storage provider
// // 	// Initialize repositories
// // 	// Initialize services

// // 	// Initialize handlers
// // 	// Setup routes

// // 	// 404 handler
// // 	return nil, nil
// // }


// package app

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/imraushankr/bervity/server/src/configs"
// 	"github.com/imraushankr/bervity/server/src/internal/handlers/v1"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/database"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/email"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/storage"
// 	"github.com/imraushankr/bervity/server/src/internal/repository"
// 	"github.com/imraushankr/bervity/server/src/internal/routes"
// 	"github.com/imraushankr/bervity/server/src/internal/services"
// )

// type RouterDependencies struct {
// }

// func SetupRouter(cfg *configs.Config, db *database.DB, log logger.Logger) (*gin.Engine, error) {
// 	router := gin.Default()

// 	// Initialize auth service
// 	authService := auth.NewAuth(&cfg.JWT)

// 	// Initialize email service
// 	emailService := email.NewEmailService(&cfg.Email, log)

// 	// Initialize storage provider
// 	storageProvider, err := storage.NewStorage(cfg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Initialize repositories
// 	userRepo := repository.NewUserRepository(db.DB, log)

// 	// Initialize services
// 	userService := services.NewUserService(
// 		userRepo,
// 		authService,
// 		emailService,
// 		cfg,
// 		storageProvider,
// 		log,
// 	)

// 	// Initialize handlers
// 	userHandler := v1.NewUserHandler(userService, log)
// 	healthHandler := v1.NewHealthHandler(cfg)

// 	// Setup routes
// 	routes.SetupRoutes(router, userHandler, healthHandler)

// 	// 404 handler
// 	router.NoRoute(func(c *gin.Context) {
// 		c.JSON(404, gin.H{"message": "Not found"})
// 	})

// 	return router, nil
// }

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

type RouterDependencies struct {
}

func SetupRouter(cfg *configs.Config, db *database.DB, log logger.Logger) (*gin.Engine, error) {
	router := gin.Default()

	// Initialize auth service
	authService := auth.NewAuth(&cfg.JWT)

	// Initialize email service
	emailService := email.NewEmailService(&cfg.Email, log)

	// Initialize storage provider
	storageProvider, err := storage.NewStorage(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB, log)

	// Initialize services
	userService := services.NewUserService(
		userRepo,
		authService,
		emailService,
		cfg,
		storageProvider,
		log,
	)

	// Initialize handlers
	userHandler := v1.NewUserHandler(userService, log)
	healthHandler := v1.NewHealthHandler(cfg)

	// Setup routes with auth dependencies
	routes.SetupRoutes(router, userHandler, healthHandler, authService, cfg)

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	return router, nil
}