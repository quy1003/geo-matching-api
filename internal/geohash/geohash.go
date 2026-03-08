package geohash

import "github.com/mmcloughlin/geohash"

const DefaultPrecision = 7

// Encode encodes a latitude and longitude into a geohash string with default precision.
func Encode(lat, lng float64) string {
	return geohash.EncodeWithPrecision(lat, lng, DefaultPrecision)
}

// EncodeWithPrecision encodes a latitude and longitude into a geohash string with the given precision.
func EncodeWithPrecision(lat, lng float64, precision uint) string {
	return geohash.EncodeWithPrecision(lat, lng, precision)
}

// Neighbors returns the 8 neighboring geohash cells for the given geohash string.
func Neighbors(hash string) []string {
	return geohash.Neighbors(hash)
}

// NeighborsWithSelf returns the given geohash cell and its 8 neighbors (total 9 cells).
func NeighborsWithSelf(hash string) []string {
	return append(geohash.Neighbors(hash), hash)
}
