package driver

// Driver represents a driver with their current location.
type Driver struct {
	ID      string  `json:"driverId"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Geohash string  `json:"geohash"`
}

// LocationUpdateRequest is the request body for updating a driver's location.
type LocationUpdateRequest struct {
	DriverID string  `json:"driverId" binding:"required"`
	Lat      float64 `json:"lat" binding:"required"`
	Lng      float64 `json:"lng" binding:"required"`
}
