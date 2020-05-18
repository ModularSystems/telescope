package daemon

import (
	"fmt"
	"log"

	"github.com/modularsystems/telescope/pkg/scan"
)

// Storage provides a high level abstraction so we can plug different types of storage (in memory, on disk, TSDB)
type Storage interface {
	Last(string) (scan.Scanner, error)
	Save(scan.Scanner)
	Size() int
	SizeOf(string) int
}

// InMemoryStore is just holding scans in memory
type InMemoryStore struct {
	CacheLength int
	Debug       bool
	Logger      *log.Logger
	ScanCache   map[string][]scan.Scanner
}

// Last returns the last scan in the map for a given key
func (i *InMemoryStore) Last(k string) (scan.Scanner, error) {
	last := len(i.ScanCache[k]) - 1
	if last < 0 {
		return &scan.HTMLScan{}, fmt.Errorf("Scan Cache is empty for %s\t", k)
	}
	if i.Debug {
		i.Logger.Printf("Found last scan for %s which occured at %s\n", k, i.ScanCache[k][last].GetTimestamp())
	}
	return i.ScanCache[k][last], nil
}

// Save manages the ScanCache, to make sure we're cycling data out
func (i *InMemoryStore) Save(s scan.Scanner) {
	// initialize cache if needed
	if len(i.ScanCache) == 0 {
		i.ScanCache = make(map[string][]scan.Scanner)
	}

	scanURI := s.GetURI()
	i.ScanCache[scanURI] = append(i.ScanCache[scanURI], s)

	// If we're over our cache size, prune the cache
	cacheSize := len(i.ScanCache[scanURI])
	if cacheSize > i.CacheLength {
		i.ScanCache[scanURI] = i.ScanCache[scanURI][1:cacheSize]
	}
	if i.Debug {
		i.Logger.Printf("Saved scan to i.ScanCache[%s]\n", scanURI)
	}
}

// Size helps to know when to initialize the map
func (i *InMemoryStore) Size() int {
	return len(i.ScanCache)
}

// SizeOf allows us to compare a key's scan slice size
func (i *InMemoryStore) SizeOf(scanURI string) int {
	return len(i.ScanCache[scanURI])
}
