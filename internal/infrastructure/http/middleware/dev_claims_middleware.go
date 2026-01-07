package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("DEV_BYPASS_AUTH") == "1" {
			_, hasTenant := c.Get("tenant_id")
			_, hasUserID := c.Get("user_id")
			_, hasUserName := c.Get("user_name")

			if !hasTenant {
				c.Set("tenant_id", "dev-tenant")
			}
			if !hasUserID {
				c.Set("user_id", "dev-user")
			}
			if !hasUserName {
				c.Set("user_name", "dev")
			}
		}

		c.Next()
	}
}
