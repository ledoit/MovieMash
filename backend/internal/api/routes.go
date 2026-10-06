package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
		} else if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// SetupRouter configures the Gin engine and API routes.
func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), cors())

	v1 := router.Group("/api/v1")
	v1.GET("/health", Health)
	v1.GET("/comparison", GetComparison)
	v1.POST("/votes", CreateVote)
	v1.GET("/leaderboard/top4", GetTop4Leaderboard)
	v1.GET("/leaderboard/movies", GetMoviesLeaderboard)

	return router
}
