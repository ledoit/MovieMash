package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ledoit/proji/backend/internal/cache"
)

type Handler struct {
	db       *sql.DB
	cache    *cache.ComparisonCache
}

func SetupRoutes(router *gin.Engine, db *sql.DB, comparisonCache *cache.ComparisonCache) {
	h := &Handler{
		db:    db,
		cache: comparisonCache,
	}

	api := router.Group("/api/v1")
	{
		api.GET("/comparison", h.GetComparison)
		api.POST("/votes", h.CreateVote)
		api.GET("/recommendations/:user_id", h.GetRecommendations)
		api.GET("/leaderboard", h.GetLeaderboard)
		api.GET("/movies/:id", h.GetMovie)
		api.POST("/auth/register", h.Register)
		api.POST("/auth/login", h.Login)
	}
}

