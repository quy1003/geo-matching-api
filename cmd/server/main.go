package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/quy1003/geo-matching-api/internal/driver"
	"github.com/quy1003/geo-matching-api/internal/handler"
	"github.com/quy1003/geo-matching-api/internal/matching"
)

func main() {
	// Wire up dependencies.
	repo := driver.NewInMemoryRepository()
	driverSvc := driver.NewService(repo)
	matcher := matching.NewMatcher(driverSvc)
	h := handler.New(driverSvc, matcher)

	r := gin.Default()

	// CORS – allow the Next.js frontend running on localhost:3000.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	h.RegisterRoutes(r)

	log.Println("Starting geo-matching-api on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
