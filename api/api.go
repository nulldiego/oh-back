package api

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nulldiego/oh-back/internal/handlers"
	"github.com/nulldiego/oh-back/internal/middleware"
	"github.com/nulldiego/oh-back/internal/model"
)

func SetupApi() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowOrigins = []string{"http://192.168.1.115:3000"}
	r.Use(cors.New(config))

	// Ping test
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})

	userRoutes := r.Group("/api/user")
	{
		userRoutes.POST("/signup", handlers.Signup)
		userRoutes.POST("/login", handlers.Login)
	}

	// Chat routes
	chatRoutes := r.Group("/api/chats", middleware.RequireAuth)
	{
		// Authorized test
		chatRoutes.GET("test", func(c *gin.Context) {
			user, exists := c.Get("user")
			if !exists {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello, %s!", user.(*model.User).Username)})
		})

		chatRoutes.GET("/", handlers.GetUserChats)        // Get current user chats
		chatRoutes.GET("/:chatId", handlers.GetMessages)  // Get chat messages
		chatRoutes.POST("/message", handlers.SendMessage) // Send message (optional chatid, if empty create new chat)
	}

	return r
}
