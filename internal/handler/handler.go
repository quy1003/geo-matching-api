package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quy1003/geo-matching-api/internal/service"
)

type Handler struct {
	datasetService  service.DatasetService
	matchingService service.MatchingService
}

func NewHandler(datasetService service.DatasetService, matchingService service.MatchingService) *Handler {
	return &Handler{
		datasetService:  datasetService,
		matchingService: matchingService,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	datasetGroup := r.Group("/datasets")
	{
		datasetGroup.POST("/upload", h.UploadDataset)
		datasetGroup.GET("", h.ListDatasets)
		datasetGroup.GET("/:id", h.GetDataset)
		datasetGroup.DELETE("/:id", h.DeleteDataset)
	}

	geoMatchingGroup := r.Group("/geo-matching")
	{
		geoMatchingGroup.POST("/run", h.RunMatching)
	}
}
