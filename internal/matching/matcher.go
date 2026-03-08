package matching

import (
	"math"

	"github.com/quy1003/geo-matching-api/internal/driver"
	gh "github.com/quy1003/geo-matching-api/internal/geohash"
)

const (
	earthRadiusKm = 6371.0
	// searchPrecision controls the geohash cell size used to build the search
	// area around the passenger.  Precision 5 gives cells of roughly 5×5 km,
	// which is a reasonable initial candidate window for ride-hailing.
	searchPrecision uint = 5
)

// RideRequest holds a passenger's location for a ride request.
type RideRequest struct {
	Lat float64 `json:"lat" binding:"required"`
	Lng float64 `json:"lng" binding:"required"`
}

// MatchResult is returned when a nearest driver is found.
type MatchResult struct {
	Driver     *driver.Driver `json:"driver"`
	DistanceKm float64        `json:"distanceKm"`
}

// Matcher encapsulates the driver matching logic.
type Matcher struct {
	driverSvc *driver.Service
}

// NewMatcher creates a new Matcher.
func NewMatcher(driverSvc *driver.Service) *Matcher {
	return &Matcher{driverSvc: driverSvc}
}

// FindNearest finds the nearest driver to the given passenger coordinates.
// It performs a geohash prefix search across the passenger's cell and its 8
// neighbors, then ranks candidates by exact Haversine distance.
func (m *Matcher) FindNearest(passengerLat, passengerLng float64) (*MatchResult, bool) {
	passengerHash := gh.EncodeWithPrecision(passengerLat, passengerLng, searchPrecision)
	searchCells := gh.NeighborsWithSelf(passengerHash)

	cellSet := make(map[string]struct{}, len(searchCells))
	for _, c := range searchCells {
		cellSet[c] = struct{}{}
	}

	allDrivers := m.driverSvc.GetAll()

	var nearest *driver.Driver
	minDist := math.MaxFloat64

	for _, d := range allDrivers {
		if !inSearchArea(d.Geohash, cellSet) {
			continue
		}
		dist := haversine(passengerLat, passengerLng, d.Lat, d.Lng)
		if dist < minDist {
			minDist = dist
			nearest = d
		}
	}

	if nearest == nil {
		return nil, false
	}
	return &MatchResult{Driver: nearest, DistanceKm: minDist}, true
}

// inSearchArea checks whether a driver geohash falls within the search cells.
// It checks both exact and prefix matches to handle precision differences.
func inSearchArea(driverHash string, cellSet map[string]struct{}) bool {
	// Exact match at the stored precision.
	if _, ok := cellSet[driverHash]; ok {
		return true
	}
	// Allow a driver hash that is a longer version of a search cell prefix.
	for cell := range cellSet {
		if len(driverHash) >= len(cell) && driverHash[:len(cell)] == cell {
			return true
		}
	}
	return false
}

// haversine computes the great-circle distance in kilometres between two
// geographic coordinates.
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	dLat := toRad(lat2 - lat1)
	dLng := toRad(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

func toRad(deg float64) float64 {
	return deg * math.Pi / 180
}
