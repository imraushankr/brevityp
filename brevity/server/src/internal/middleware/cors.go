// package middleware

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/imraushankr/bervity/server/src/configs"
// )

// func CORSMiddleware(cfg *configs.Config) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.CORS.AllowOrigins[])
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }

package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
)

func CORSMiddleware(cfg *configs.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.CORS.Enabled {
			c.Next()
			return
		}

		// Handle multiple allowed origins
		origin := c.Request.Header.Get("Origin")
		allowedOrigin := ""
		for _, o := range cfg.CORS.AllowOrigins {
			if o == "*" || o == origin {
				allowedOrigin = o
				break
			}
		}

		if allowedOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join([]string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"accept",
			"origin",
			"Cache-Control",
			"X-Requested-With",
		}, ","))
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.CORS.AllowMethods, ","))
		c.Writer.Header().Set("Access-Control-Max-Age", cfg.CORS.MaxAge)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
