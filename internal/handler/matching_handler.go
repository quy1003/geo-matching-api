package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quy1003/geo-matching-api/internal/dto"
	"github.com/quy1003/geo-matching-api/internal/service"
)

func (h *Handler) RunMatching(c *gin.Context) {
	var req dto.RunMatchingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.matchingService.Run(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRadius):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrDatasetNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "dataset not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "matching failed"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}
