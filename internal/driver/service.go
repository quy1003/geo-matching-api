package driver

import gh "github.com/quy1003/geo-matching-api/internal/geohash"

// Service handles business logic for driver operations.
type Service struct {
	repo Repository
}

// NewService creates a new driver Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// UpdateLocation encodes the driver's coordinates into a geohash and persists the record.
func (s *Service) UpdateLocation(id string, lat, lng float64) *Driver {
	hash := gh.Encode(lat, lng)
	d := &Driver{
		ID:      id,
		Lat:     lat,
		Lng:     lng,
		Geohash: hash,
	}
	s.repo.Save(d)
	return d
}

// GetAll returns all drivers currently stored.
func (s *Service) GetAll() []*Driver {
	return s.repo.FindAll()
}
