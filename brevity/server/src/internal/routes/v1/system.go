package v1

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/imraushankr/bervity/server/src/internal/handlers/v1"
	"github.com/imraushankr/bervity/server/src/internal/middleware"
)

func RegisterSystemRoutes(r *gin.RouterGroup, h *v1.HealthHandler) {
	sys := r.Group("/system")
	{
		// Health checks
		sys.GET("/health", h.HealthCheck)
		sys.GET("/status", h.GetStatus)

		// Metrics and monitoring
		sys.GET("/metrics", middleware.PrometheusMetricsMiddleware(), middleware.PrometheusHandler())
		sys.GET("/stats", h.GetStatistics)

		// Configuration (protected)
		sys.GET("/config", h.GetConfig)
	}
}