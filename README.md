# geo-matching-api

A backend service written in **Go** that matches a passenger with the nearest available driver using **Geohash-based spatial indexing** — inspired by simplified ride-hailing dispatch systems.

---

## System Purpose

The API simulates two core operations of a ride-hailing platform:

1. **Driver location updates** – drivers periodically report their current GPS coordinates; the service encodes each position as a geohash and stores it in memory.
2. **Passenger ride requests** – a passenger submits their GPS coordinates; the service identifies candidate drivers via geohash neighbor search, ranks them by exact Haversine distance, and returns the nearest one.

---

## How Geohash Matching Works

[Geohash](https://en.wikipedia.org/wiki/Geohash) is a hierarchical spatial index that encodes a latitude/longitude pair into a short alphanumeric string. Adjacent cells share a common prefix, making prefix-based proximity queries fast and simple.

### Matching algorithm

```
Passenger (lat, lng)
       │
       ▼
Encode to geohash at precision 5 (~5 km × 5 km cell)
       │
       ▼
Collect the cell + 8 neighbouring cells  ←  3×3 search grid ≈ 225 km²
       │
       ▼
Filter drivers whose stored geohash starts with any search-cell prefix
       │
       ▼
Compute exact Haversine distance for each candidate
       │
       ▼
Return the driver with the smallest distance
```

**Why precision 5 for the search area?**
A precision-5 cell covers roughly 4.9 km × 4.9 km. Using a 3×3 neighbourhood therefore covers a ~15 km × 15 km window — a sensible initial candidate window for urban ride-hailing.

Driver locations are stored at **precision 7** (~76 m × 38 m) for accurate position records; the search simply checks whether a driver's hash begins with the lower-precision search prefix.

---

## Project Structure

```
geo-matching-api/
├── cmd/
│   └── server/
│       └── main.go          # HTTP server entry-point
├── internal/
│   ├── driver/
│   │   ├── model.go         # Driver struct & request types
│   │   ├── repository.go    # Thread-safe in-memory store
│   │   └── service.go       # Business logic (encode + persist)
│   ├── geohash/
│   │   └── geohash.go       # Encode / neighbor-lookup wrappers
│   ├── matching/
│   │   └── matcher.go       # Geohash search + Haversine ranking
│   └── handler/
│       └── handler.go       # Gin HTTP handlers
└── go.mod
```

---

## API Reference

### `GET /health`

Returns service health status.

**Response**
```json
{ "status": "ok" }
```

---

### `POST /drivers/location`

Updates a driver's current location.

**Request body**
```json
{
  "driverId": "driver-42",
  "lat": 10.7769,
  "lng": 106.7009
}
```

**Response** – the stored driver record including the computed geohash:
```json
{
  "driverId": "driver-42",
  "lat": 10.7769,
  "lng": 106.7009,
  "geohash": "w3gv6kz"
}
```

---

### `POST /rides/request`

Requests a ride for a passenger at the given coordinates.

**Request body**
```json
{
  "lat": 10.7769,
  "lng": 106.7009
}
```

**Response** – the nearest driver and their distance:
```json
{
  "driver": {
    "driverId": "driver-42",
    "lat": 10.7800,
    "lng": 106.7020,
    "geohash": "w3gv6mz"
  },
  "distanceKm": 0.42
}
```

Returns **404** when no drivers are available in the search area.

---

## Running the Server Locally

### Prerequisites

* [Go 1.21+](https://golang.org/dl/)

### Steps

```bash
# 1. Clone the repository
git clone https://github.com/quy1003/geo-matching-api.git
cd geo-matching-api

# 2. Download dependencies
go mod download

# 3. Start the server (listens on :8080)
go run ./cmd/server/main.go
```

The server starts on **http://localhost:8080**.

### Example requests

```bash
# Health check
curl http://localhost:8080/health

# Register a driver
curl -X POST http://localhost:8080/drivers/location \
  -H "Content-Type: application/json" \
  -d '{"driverId":"driver-1","lat":10.7800,"lng":106.7020}'

# Request a ride
curl -X POST http://localhost:8080/rides/request \
  -H "Content-Type: application/json" \
  -d '{"lat":10.7769,"lng":106.7009}'
```

### Running tests

```bash
go test ./...
```

---

## CORS

The server allows cross-origin requests from `http://localhost:3000` (Next.js development server) for all standard HTTP methods and headers.

---

## Tech Stack

| Concern | Library |
|---|---|
| HTTP framework | [gin-gonic/gin](https://github.com/gin-gonic/gin) |
| CORS middleware | [gin-contrib/cors](https://github.com/gin-contrib/cors) |
| Geohash encoding | [mmcloughlin/geohash](https://github.com/mmcloughlin/geohash) |
| Storage | In-memory (thread-safe `sync.RWMutex` map) |

---

## Notes

* Driver data is stored **in-memory** and will be lost on server restart. A future version could swap the repository layer for Redis, PostgreSQL with PostGIS, or another persistent store without changing the matching or handler layers.
* The matching search radius (precision 5 ≈ 15 km window) is intentionally conservative. It can be tuned via the `searchPrecision` constant in `internal/matching/matcher.go`.

