package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

// JWTAuth creates a Gin middleware for JWT authentication
func JWTAuth(authService *auth.Auth, cfg *configs.Config, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header or cookie
		tokenString := extractToken(c, cfg.JWT.SecureCookie)
		if tokenString == "" {
			log.Warn("No authentication token provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Authorization token required",
			})
			return
		}

		// Validate token
		claims, err := authService.VerifyAccessToken(tokenString)
		if err != nil {
			status := http.StatusUnauthorized
			if errors.Is(err, models.ErrExpiredToken) {
				status = http.StatusForbidden
				if claims != nil {
					log.Warn("Expired token attempt", logger.String("user_id", claims.UserId))
				}
			} else {
				log.Warn("Invalid token attempt", logger.NamedError("error", err))
			}

			c.AbortWithStatusJSON(status, models.ErrorResponse{
				Error: "Invalid token",
			})
			return
		}

		// Ensure claims are valid
		if claims == nil || claims.UserId == "" {
			log.Warn("Token validation returned empty claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid token claims",
			})
			return
		}

		// Set user context
		c.Set("user_id", claims.UserId)
		c.Set("user_role", claims.Role)
		log.Debug("User authenticated", 
			logger.String("user_id", claims.UserId),
			logger.String("role", claims.Role))

		c.Next()
	}
}

// RoleAuth creates a middleware to check user roles
func RoleAuth(allowedRoles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get logger from Gin context if available, otherwise create new one
		var log logger.Logger
		if val, exists := c.Get("logger"); exists {
			if loggerVal, ok := val.(logger.Logger); ok {
				log = loggerVal
			}
		}
		if log == nil {
			log = logger.Get()
		}

		userRole, exists := c.Get("user_role")
		if !exists {
			log.Warn("Role check failed - no role information")
			c.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse{
				Error: "Forbidden - no role information",
			})
			return
		}

		// Convert to Role type
		role, ok := userRole.(models.Role)
		if !ok {
			// Try string conversion if direct type assertion fails
			if roleStr, isString := userRole.(string); isString {
				role = models.Role(roleStr)
			} else {
				log.Error("Invalid role type in context")
				c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{
					Error: "Internal server error",
				})
				return
			}
		}

		// Check if user has any of the allowed roles
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				log.Debug("Role authorized", logger.String("role", string(role)))
				c.Next()
				return
			}
		}

		log.Warn("Insufficient permissions", 
			logger.String("required_roles", strings.Join(rolesToStrings(allowedRoles), ", ")),
			logger.String("user_role", string(role)))

		c.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse{
			Error: "Forbidden - insufficient permissions",
		})
	}
}

// RefreshTokenAuth creates a Gin middleware for refresh token authentication
func RefreshTokenAuth(authService *auth.Auth, cfg *configs.Config, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie (refresh tokens should only come via secure cookie)
		tokenString, err := c.Cookie("refresh_token")
		if err != nil {
			log.Warn("Refresh token missing from cookies")
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Refresh token required",
			})
			return
		}

		// Validate refresh token
		claims, err := authService.VerifyRefreshToken(tokenString)
		if err != nil {
			status := http.StatusUnauthorized
			if errors.Is(err, models.ErrExpiredToken) {
				status = http.StatusForbidden
				if claims != nil {
					log.Warn("Expired refresh token attempt", logger.String("user_id", claims.UserId))
				}
			} else {
				log.Warn("Invalid refresh token attempt", logger.NamedError("error", err))
			}

			c.AbortWithStatusJSON(status, models.ErrorResponse{
				Error: "Invalid refresh token",
			})
			return
		}

		// Ensure claims are valid
		if claims == nil || claims.UserId == "" {
			log.Warn("Refresh token validation returned empty claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid refresh token claims",
			})
			return
		}

		// Set user context
		c.Set("user_id", claims.UserId)
		c.Set("user_role", claims.Role)
		log.Debug("Refresh token validated", 
			logger.String("user_id", claims.UserId),
			logger.String("role", claims.Role))

		c.Next()
	}
}

// extractToken tries to get the JWT token from different sources
func extractToken(c *gin.Context, secureCookie bool) string {
	// 1. Check Authorization header first
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	// 2. Fallback to cookie if enabled
	if secureCookie {
		if token, err := c.Cookie("access_token"); err == nil && token != "" {
			return token
		}
	}

	// 3. Check query parameter (for websockets or special cases)
	if token := c.Query("token"); token != "" {
		return token
	}

	return ""
}

// rolesToStrings converts a slice of Role to a slice of strings
func rolesToStrings(roles []models.Role) []string {
	var strRoles []string
	for _, role := range roles {
		strRoles = append(strRoles, string(role))
	}
	return strRoles
}