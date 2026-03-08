package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quy1003/geo-matching-api/internal/driver"
	"github.com/quy1003/geo-matching-api/internal/handler"
	"github.com/quy1003/geo-matching-api/internal/matching"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	repo := driver.NewInMemoryRepository()
	driverSvc := driver.NewService(repo)
	matcher := matching.NewMatcher(driverSvc)
	h := handler.New(driverSvc, matcher)

	r := gin.New()
	h.RegisterRoutes(r)
	return r
}

func TestHealthCheck(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %s", body["status"])
	}
}

func TestUpdateDriverLocation_Valid(t *testing.T) {
	r := setupRouter()
	payload := `{"driverId":"d1","lat":10.7769,"lng":106.7009}`
	req := httptest.NewRequest(http.MethodPost, "/drivers/location", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var d driver.Driver
	if err := json.NewDecoder(w.Body).Decode(&d); err != nil {
		t.Fatalf("failed to decode driver: %v", err)
	}
	if d.ID != "d1" {
		t.Fatalf("expected driverId d1, got %s", d.ID)
	}
	if d.Geohash == "" {
		t.Fatal("expected non-empty geohash")
	}
}

func TestUpdateDriverLocation_MissingField(t *testing.T) {
	r := setupRouter()
	payload := `{"lat":10.7769,"lng":106.7009}`
	req := httptest.NewRequest(http.MethodPost, "/drivers/location", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing driverId, got %d", w.Code)
	}
}

func TestRequestRide_NoDrivers(t *testing.T) {
	r := setupRouter()
	payload := `{"lat":10.7769,"lng":106.7009}`
	req := httptest.NewRequest(http.MethodPost, "/rides/request", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 when no drivers, got %d", w.Code)
	}
}

func TestRequestRide_WithDriver(t *testing.T) {
	r := setupRouter()

	// Register a driver first.
	locPayload := `{"driverId":"d1","lat":10.7769,"lng":106.7009}`
	locReq := httptest.NewRequest(http.MethodPost, "/drivers/location", bytes.NewBufferString(locPayload))
	locReq.Header.Set("Content-Type", "application/json")
	locW := httptest.NewRecorder()
	r.ServeHTTP(locW, locReq)
	if locW.Code != http.StatusOK {
		t.Fatalf("setup: failed to register driver: %d", locW.Code)
	}

	// Now request a ride from the same location.
	ridePayload := `{"lat":10.7769,"lng":106.7009}`
	rideReq := httptest.NewRequest(http.MethodPost, "/rides/request", bytes.NewBufferString(ridePayload))
	rideReq.Header.Set("Content-Type", "application/json")
	rideW := httptest.NewRecorder()
	r.ServeHTTP(rideW, rideReq)

	if rideW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rideW.Code, rideW.Body.String())
	}

	var result matching.MatchResult
	if err := json.NewDecoder(rideW.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode match result: %v", err)
	}
	if result.Driver == nil {
		t.Fatal("expected a driver in the result")
	}
	if result.Driver.ID != "d1" {
		t.Fatalf("expected driver d1, got %s", result.Driver.ID)
	}
}

func TestRequestRide_InvalidBody(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodPost, "/rides/request", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
