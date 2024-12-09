package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nulldiego/oh-back/internal/middleware"
)

func GetCurrentUser(c *gin.Context) *middleware.AuthUser {
	cUser, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current user"})
		return nil
	}

	if user, ok := cUser.(middleware.AuthUser); ok {
		return &user
	}

	return nil
}
