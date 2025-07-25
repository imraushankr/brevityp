// middleware/anonymous_url.go
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

func AnonymousURLLimit(urlRepo interfaces.URLRepository, log logger.Logger, limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if user is authenticated
		if userID, exists := c.Get("user_id"); exists && userID != "" {
			c.Next()
			return
		}

		// Get client IP
		clientIP := c.ClientIP()

		// Count URLs created by this IP
		count, err := urlRepo.CountByIP(c.Request.Context(), clientIP)
		if err != nil {
			log.Error("failed to count anonymous URLs", logger.ErrorField(err))
			c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
			return
		}

		if count >= limit {
			c.AbortWithStatusJSON(403, gin.H{
				"error":   "anonymous_url_limit_reached",
				"message": "You've reached the limit of 5 anonymous URLs. Please login to create more.",
			})
			return
		}

		c.Next()
	}
}
