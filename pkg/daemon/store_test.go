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
	}
	store.Save(scan)
	assert.Equal(t, len(i.ScanCache), 1, "InMemoryStore failed to store the submitted scan")

	lastScanInMemory := i.ScanCache[scan.ScanName][0]
	assert.Equal(t, scan.ScanName, lastScanInMemory.Name(), "The saved scan has a different name than what was last saved.")

	lastScan := store.Last(scan)
	assert.Equal(t, lastScan.Name(), scan.ScanName, "store didn't seem to have our last scan")
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
		}
		store.Save(scan)
	}
	assert.Equal(t, len(i.ScanCache["testscan"]), 10, "InMemoryStore stored an incorrect number of entries, expected 10, got %d", len(i.ScanCache["testscan"]))
}
