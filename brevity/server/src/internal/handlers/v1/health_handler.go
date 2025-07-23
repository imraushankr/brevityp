package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
)

type HealthHandler struct {
	cfg *configs.Config
}

func NewHealthHandler(cfg *configs.Config) *HealthHandler {
	return &HealthHandler{cfg: cfg}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"system": gin.H{
			"version":     h.cfg.App.Version,
			"environment": h.cfg.App.Environment,
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func (h *HealthHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"app": gin.H{
			"name":        h.cfg.App.Name,
			"version":     h.cfg.App.Version,
			"environment": h.cfg.App.Environment,
			"debug":       h.cfg.App.Debug,
		},
		"server": gin.H{
			"host":             h.cfg.Server.Host,
			"port":             h.cfg.Server.Port,
			"read_timeout":     h.cfg.Server.ReadTimeout.String(),
			"write_timeout":    h.cfg.Server.WriteTimeout.String(),
			"shutdown_timeout": h.cfg.Server.ShutdownTimeout.String(),
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *HealthHandler) GetMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"metrics": "Not implemented yet",
	})
}

func (h *HealthHandler) GetStatistics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"stats": "Not implemented yet",
	})
}

func (h *HealthHandler) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, h.cfg)
}