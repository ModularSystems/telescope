package daemon

import (
	"testing"

	"github.com/modularsystems/telescope/pkg/scan"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryStore(t *testing.T) {
	var store Storage
	i := &InMemoryStore{
		CacheLength: 100,
	}
	store = i
	assert.NotNil(t, store)

	scan := &scan.HTMLScan{
		ScanName: "testscan",
		URI:      "localhost",
	}
	store.Save(scan)
	assert.Equal(t, len(i.ScanCache), 1, "InMemoryStore failed to store the submitted scan")

	lastScanInMemory := i.ScanCache[scan.URI][0]
	assert.Equal(t, scan.ScanName, lastScanInMemory.GetName(), "The saved scan has a different name than what was last saved.")

	lastScan, err := store.Last(scan.URI)
	assert.NoErrorf(t, err, "Expected to get no errors when calling store.Last")
	assert.Equal(t, lastScan.GetURI(), scan.URI, "store didn't seem to have our last scan")
}

// TestInMemoryStoreCachePurging creates an in memory store that holds 10 scans, then saves 11 scans, and
// ensures the last 10 were saved
func TestInMemoryStoreCachePurging(t *testing.T) {
	var store Storage
	i := &InMemoryStore{
		CacheLength: 10,
	}
	store = i
	assert.NotNil(t, store)

	for j := 0; j < 100; j++ {
		scan := &scan.HTMLScan{
			ScanName: "testscan",
			URI:      "localhost",
		}
		store.Save(scan)
	}
	assert.Equal(t, len(i.ScanCache["localhost"]), 10, "InMemoryStore stored an incorrect number of entries, expected 10, got %d", len(i.ScanCache["testscan"]))
}
