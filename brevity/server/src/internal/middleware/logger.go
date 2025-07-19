package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

func Logs(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		log.Info("HTTP Request",
			logger.Int("status", c.Writer.Status()),
			logger.String("method", c.Request.Method),
			logger.String("path", path),
			logger.String("query", query),
			logger.String("ip", c.ClientIP()),
			logger.String("user-agent", c.Request.UserAgent()),
		)
	}
}
