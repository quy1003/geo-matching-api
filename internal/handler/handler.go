package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quy1003/geo-matching-api/internal/driver"
	"github.com/quy1003/geo-matching-api/internal/matching"
)

// Handler holds all HTTP handler dependencies.
type Handler struct {
	driverSvc *driver.Service
	matcher   *matching.Matcher
}

// New creates a new Handler.
func New(driverSvc *driver.Service, matcher *matching.Matcher) *Handler {
	return &Handler{driverSvc: driverSvc, matcher: matcher}
}

// RegisterRoutes attaches all API routes to the given Gin engine.
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.HealthCheck)
	r.POST("/drivers/location", h.UpdateDriverLocation)
	r.POST("/rides/request", h.RequestRide)
}

// HealthCheck godoc
// @Summary Health check
// @Description Returns the service health status
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// UpdateDriverLocation godoc
// @Summary Update driver location
// @Description Stores or updates a driver's current geohash-encoded location
// @Accept json
// @Produce json
// @Param body body driver.LocationUpdateRequest true "Driver location"
// @Success 200 {object} driver.Driver
// @Failure 400 {object} map[string]string
// @Router /drivers/location [post]
func (h *Handler) UpdateDriverLocation(c *gin.Context) {
	var req driver.LocationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	d := h.driverSvc.UpdateLocation(req.DriverID, req.Lat, req.Lng)
	c.JSON(http.StatusOK, d)
}

// RequestRide godoc
// @Summary Request a ride
// @Description Finds the nearest available driver using geohash spatial search
// @Accept json
// @Produce json
// @Param body body matching.RideRequest true "Passenger location"
// @Success 200 {object} matching.MatchResult
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /rides/request [post]
func (h *Handler) RequestRide(c *gin.Context) {
	var req matching.RideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, found := h.matcher.FindNearest(req.Lat, req.Lng)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "no drivers available nearby"})
		return
	}

	c.JSON(http.StatusOK, result)
}
