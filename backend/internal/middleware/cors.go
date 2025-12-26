package middleware

import (
	"gin-quickstart/internal/config"

	"github.com/gin-gonic/gin"
)

// CORS creates a middleware for handling CORS
func CORS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		if len(cfg.AllowedOrigins) == 0 {
			// If no origins configured, allow all (development mode)
			allowed = true
			origin = "*"
		} else {
			for _, allowedOrigin := range cfg.AllowedOrigins {
				if allowedOrigin == origin || allowedOrigin == "*" {
					allowed = true
					break
				}
			}
		}

		if allowed {
			// Set CORS headers
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
