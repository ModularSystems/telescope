package daemon

import (
	"github.com/modularsystems/telescope/pkg/scan"
)

// Storage provides a high level abstraction so we can plug different types of storage (in memory, on disk, TSDB)
type Storage interface {
	Last(scan.Scanner) scan.Scanner
	Save(scan.Scanner)
	Size() int
	SizeOf(string) int
}

// InMemoryStore is just holding scans in memory
type InMemoryStore struct {
	CacheLength int
	ScanCache   map[string][]scan.Scanner
}

// Last returns the last scan in the map
func (i *InMemoryStore) Last(s scan.Scanner) scan.Scanner {
	last := len(i.ScanCache[s.GetName()]) - 1
	return i.ScanCache[s.GetName()][last]
}

// Save manages the ScanCache, to make sure we're cycling data out
func (i *InMemoryStore) Save(s scan.Scanner) {
	// initialize cache if needed
	if len(i.ScanCache) == 0 {
		i.ScanCache = make(map[string][]scan.Scanner)
	}

	scanName := s.GetName()
	i.ScanCache[scanName] = append(i.ScanCache[scanName], s)

	// If we're over our cache size, prune the cache
	cacheSize := len(i.ScanCache[scanName])
	if cacheSize > i.CacheLength {
		i.ScanCache[scanName] = i.ScanCache[s.GetName()][1:cacheSize]
	}
}

func (i *InMemoryStore) Size() int {
	return len(i.ScanCache)
}

func (i *InMemoryStore) SizeOf(scanName string) int {
	return len(i.ScanCache[scanName])
}
