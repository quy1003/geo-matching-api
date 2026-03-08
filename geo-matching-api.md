# GeoHash Matching API

A backend service written in Go that simulates a simplified ride-hailing dispatch system.
The system focuses on solving the **geospatial matching problem**: finding the nearest driver for a passenger using Geohash-based spatial indexing.

This project is a technical demo inspired by ride-hailing platforms such as Grab and Uber, focusing on location encoding and driver matching algorithms.

---

## Overview

In ride-hailing systems, the core challenge is efficiently finding nearby drivers among potentially thousands of active vehicles.

Instead of calculating the distance to every driver, this project uses **Geohash spatial indexing** to narrow down the search area before computing the exact distance.

Workflow:

1. Drivers periodically update their location.
2. The server encodes the location into a Geohash string.
3. When a passenger requests a ride:
   - the passenger location is encoded into Geohash
   - nearby Geohash cells are searched
   - candidate drivers are filtered

4. The system calculates the real distance using the **Haversine formula**.
5. The nearest driver is returned.

---

## Tech Stack

Backend language

- Go

Web framework

- Gin

Geospatial logic

- Geohash encoding

Architecture

- Modular Go project structure

---

## Project Structure

cmd/server
Entry point for the HTTP API server

internal/driver
Driver models and driver location management

internal/geohash
Geohash encoding utilities and neighbor lookup

internal/matching
Driver matching algorithm

internal/handler
HTTP handlers for API endpoints

---

## API Endpoints

Health check

GET /health

Response

```
API is running
```

---

Update driver location

POST /drivers/location

Request body

```
{
  "driverId": "driver-1",
  "lat": 10.762622,
  "lng": 106.660172
}
```

Behavior

- Encodes the location into Geohash
- Stores the driver location in memory

---

Request ride

POST /rides/request

Request body

```
{
  "lat": 10.762622,
  "lng": 106.660172
}
```

Behavior

- Encodes passenger location into Geohash
- Searches nearby drivers using Geohash prefix
- Calculates real distance using Haversine formula
- Returns the nearest driver

---

## Local Development

Clone the repository

```
git clone https://github.com/yourusername/geo-matching-api.git
cd geo-matching-api
```

Install dependencies

```
go mod tidy
```

Run the server

```
go run cmd/server/main.go
```

Server will start at

```
http://localhost:8080
```

---

## CORS Configuration

The API enables CORS for the frontend application running at

```
http://localhost:3000
```

This allows the Next.js frontend to call the backend API during development.

---

## Future Improvements

Possible improvements for this project

- Persist driver locations using Redis
- Implement real-time driver updates via WebSocket
- Add spatial indexing for large-scale datasets
- Simulate thousands of drivers
- Visua
