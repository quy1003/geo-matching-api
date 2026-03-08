package driver_test

import (
	"testing"

	"github.com/quy1003/geo-matching-api/internal/driver"
)

func TestInMemoryRepository(t *testing.T) {
	repo := driver.NewInMemoryRepository()

	d := &driver.Driver{ID: "d1", Lat: 10.0, Lng: 106.0, Geohash: "w3"}
	repo.Save(d)

	all := repo.FindAll()
	if len(all) != 1 {
		t.Fatalf("expected 1 driver, got %d", len(all))
	}

	got, ok := repo.FindByID("d1")
	if !ok {
		t.Fatal("expected to find driver d1")
	}
	if got.ID != "d1" {
		t.Fatalf("expected driver ID d1, got %s", got.ID)
	}

	_, ok = repo.FindByID("unknown")
	if ok {
		t.Fatal("expected no driver for unknown ID")
	}
}

func TestServiceUpdateLocation(t *testing.T) {
	repo := driver.NewInMemoryRepository()
	svc := driver.NewService(repo)

	d := svc.UpdateLocation("d1", 10.7769, 106.7009)
	if d.ID != "d1" {
		t.Fatalf("unexpected ID %s", d.ID)
	}
	if d.Geohash == "" {
		t.Fatal("expected non-empty geohash after UpdateLocation")
	}

	all := svc.GetAll()
	if len(all) != 1 {
		t.Fatalf("expected 1 driver, got %d", len(all))
	}
}
