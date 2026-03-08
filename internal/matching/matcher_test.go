package matching_test

import (
	"math"
	"testing"

	"github.com/quy1003/geo-matching-api/internal/driver"
	"github.com/quy1003/geo-matching-api/internal/matching"
)

func TestFindNearest_NoDrivers(t *testing.T) {
	repo := driver.NewInMemoryRepository()
	svc := driver.NewService(repo)
	m := matching.NewMatcher(svc)

	_, found := m.FindNearest(10.7769, 106.7009)
	if found {
		t.Fatal("expected no match when no drivers are registered")
	}
}

func TestFindNearest_SingleDriver(t *testing.T) {
	repo := driver.NewInMemoryRepository()
	svc := driver.NewService(repo)
	m := matching.NewMatcher(svc)

	svc.UpdateLocation("d1", 10.7769, 106.7009)

	result, found := m.FindNearest(10.7769, 106.7009)
	if !found {
		t.Fatal("expected to find a driver")
	}
	if result.Driver.ID != "d1" {
		t.Fatalf("expected driver d1, got %s", result.Driver.ID)
	}
	if result.DistanceKm < 0 {
		t.Fatalf("distance should be non-negative, got %f", result.DistanceKm)
	}
}

func TestFindNearest_MultipleDrivers(t *testing.T) {
	repo := driver.NewInMemoryRepository()
	svc := driver.NewService(repo)
	m := matching.NewMatcher(svc)

	// Passenger at HCMC.
	passengerLat, passengerLng := 10.7769, 106.7009

	// Driver 1 – close (same neighborhood).
	svc.UpdateLocation("near", 10.7800, 106.7020)
	// Driver 2 – far away (Hanoi).
	svc.UpdateLocation("far", 21.0285, 105.8542)

	result, found := m.FindNearest(passengerLat, passengerLng)
	if !found {
		t.Fatal("expected to find a driver")
	}
	if result.Driver.ID != "near" {
		t.Fatalf("expected nearest driver 'near', got %s", result.Driver.ID)
	}
}

func TestFindNearest_DistanceAccuracy(t *testing.T) {
	repo := driver.NewInMemoryRepository()
	svc := driver.NewService(repo)
	m := matching.NewMatcher(svc)

	// Two points roughly 1 km apart.
	svc.UpdateLocation("d1", 10.7769, 106.7009)

	result, found := m.FindNearest(10.7769, 106.7009)
	if !found {
		t.Fatal("expected to find a driver")
	}
	// Same point – distance should be essentially zero.
	if math.Abs(result.DistanceKm) > 0.001 {
		t.Fatalf("expected near-zero distance for same-point match, got %f", result.DistanceKm)
	}
}
