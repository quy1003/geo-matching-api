package app

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/quy1003/geo-matching-api/internal/handler"
	"github.com/quy1003/geo-matching-api/internal/repository/memory"
	"github.com/quy1003/geo-matching-api/internal/service"
)

func NewServer() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	store := memory.NewStore()

	datasetService := service.NewDatasetService(store)
	matchingService := service.NewMatchingService(store, store)

	h := handler.NewHandler(datasetService, matchingService)
	h.RegisterRoutes(r)

	return r
}
