package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/quy1003/geo-matching-api/internal/dto"
	"github.com/quy1003/geo-matching-api/internal/service"
)

func (h *Handler) UploadDataset(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only csv files are supported"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot read uploaded file"})
		return
	}
	defer file.Close()

	name := c.PostForm("name")
	dataset, err := h.datasetService.UploadCSV(c.Request.Context(), name, fileHeader.Filename, file)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCSVFormat), errors.Is(err, service.ErrInvalidCoordinate):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload dataset"})
		}
		return
	}

	c.JSON(http.StatusCreated, dataset)
}

func (h *Handler) ListDatasets(c *gin.Context) {
	datasets, err := h.datasetService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list datasets"})
		return
	}

	c.JSON(http.StatusOK, datasets)
}

func (h *Handler) GetDataset(c *gin.Context) {
	datasetID := c.Param("id")
	dataset, points, err := h.datasetService.Get(c.Request.Context(), datasetID)
	if err != nil {
		if errors.Is(err, service.ErrDatasetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "dataset not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get dataset"})
		return
	}

	c.JSON(http.StatusOK, dto.DatasetDetailResponse{
		Dataset: dataset,
		Points:  points,
	})
}

func (h *Handler) DeleteDataset(c *gin.Context) {
	datasetID := c.Param("id")
	err := h.datasetService.Delete(c.Request.Context(), datasetID)
	if err != nil {
		if errors.Is(err, service.ErrDatasetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "dataset not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete dataset"})
		return
	}

	c.Status(http.StatusNoContent)
}
