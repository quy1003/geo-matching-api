package geohash_test

import (
	"strings"
	"testing"

	gh "github.com/quy1003/geo-matching-api/internal/geohash"
)

func TestEncode(t *testing.T) {
	// Ho Chi Minh City, Vietnam.
	hash := gh.Encode(10.7769, 106.7009)
	if len(hash) != gh.DefaultPrecision {
		t.Fatalf("expected hash length %d, got %d", gh.DefaultPrecision, len(hash))
	}
	if !strings.HasPrefix(hash, "w3gv") {
		t.Fatalf("unexpected geohash prefix for HCMC: %s", hash)
	}
}

func TestEncodeWithPrecision(t *testing.T) {
	hash := gh.EncodeWithPrecision(10.7769, 106.7009, 5)
	if len(hash) != 5 {
		t.Fatalf("expected precision 5 hash, got length %d", len(hash))
	}
}

func TestNeighbors(t *testing.T) {
	hash := gh.Encode(10.7769, 106.7009)
	neighbors := gh.Neighbors(hash)
	if len(neighbors) != 8 {
		t.Fatalf("expected 8 neighbors, got %d", len(neighbors))
	}
}

func TestNeighborsWithSelf(t *testing.T) {
	hash := gh.Encode(10.7769, 106.7009)
	cells := gh.NeighborsWithSelf(hash)
	if len(cells) != 9 {
		t.Fatalf("expected 9 cells (self + 8 neighbors), got %d", len(cells))
	}
	found := false
	for _, c := range cells {
		if c == hash {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("NeighborsWithSelf should include the hash itself")
	}
}
