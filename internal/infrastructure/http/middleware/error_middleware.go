package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"

	"ms-parcel-core/internal/pkg/util/apperror"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() {
			return
		}

		last := c.Errors.Last()
		if last == nil {
			return
		}

		var appErr *apperror.AppError
		if errors.As(last.Err, &appErr) && appErr != nil {
			c.JSON(appErr.Status, gin.H{"success": false, "error": appErr})
			return
		}

		c.JSON(500, gin.H{
			"success": false,
			"error":   apperror.NewInternal("internal_error", "error interno", map[string]any{"error": last.Err.Error()}),
		})
	}
}
